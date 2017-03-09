#include <jni.h>
#include <stdlib.h>
#include <android/log.h>
#include <math.h>
#include <SuperpoweredFrequencyDomain.h>
#include <AndroidIO/SuperpoweredAndroidAudioIO.h>
#include <SuperpoweredSimple.h>
#include <SuperpoweredCPU.h>
#include <SLES/OpenSLES.h>
#include <SLES/OpenSLES_AndroidConfiguration.h>


static SuperpoweredFrequencyDomain *frequencyDomain;
static float *magnitudeLeft, *magnitudeRight, *phaseLeft, *phaseRight, *fifoOutput, *inputBufferFloat;
static int fifoOutputFirstSample, fifoOutputLastSample, stepSize, fifoCapacity;
JavaVM* gJvm = 0;
static jobject gClassLoader;
static jmethodID gFindClassMethod;
static jclass mClass;

#define FFT_LOG_SIZE 11 // 2^11 = 2048

JNIEnv* getEnv() {
    JNIEnv *env;
    int status = gJvm->GetEnv((void**)&env, JNI_VERSION_1_6);
    if(status < 0) {
        status = gJvm->AttachCurrentThread(&env, NULL);
        if(status < 0) {
            return 0;
        }
    }
    return env;
}



JNIEXPORT jint JNICALL JNI_OnLoad(JavaVM *pjvm, void *reserved) {
    gJvm = pjvm;  // cache the JavaVM pointer
    auto env = getEnv();
    //replace with one of your classes in the line below
    auto randomClass = env->FindClass("se/livetsord/youth/MainActivity");
    mClass = (jclass)env->NewGlobalRef(randomClass);
    jclass classClass = env->GetObjectClass(randomClass);
    auto classLoaderClass = env->FindClass("java/lang/ClassLoader");
    auto getClassLoaderMethod = env->GetMethodID(classClass, "getClassLoader",
                                                 "()Ljava/lang/ClassLoader;");
    gClassLoader = env->CallObjectMethod(randomClass, getClassLoaderMethod);
    gFindClassMethod = env->GetMethodID(classLoaderClass, "findClass",
                                        "(Ljava/lang/String;)Ljava/lang/Class;");

    return JNI_VERSION_1_6;
}

jclass findClass(const char* name) {
    return static_cast<jclass>(getEnv()->CallObjectMethod(gClassLoader, gFindClassMethod, getEnv()->NewStringUTF(name)));
}





// This is called periodically by the media server.
static bool
audioProcessing(void *__unused clientdata, short int *audioInputOutput, int numberOfSamples,
                int __unused samplerate) {
    int numberOfChannels = 24;
    int vals[numberOfChannels];
    float treashold = 0.15f;
    int rgb[3] = {0,0, 0};
    int diff = 18;
    SuperpoweredShortIntToFloat(audioInputOutput, inputBufferFloat,
                                (unsigned int) numberOfSamples); // Converting the 16-bit integer samples to 32-bit floating point.
    frequencyDomain->addInput(inputBufferFloat,
                              numberOfSamples); // Input goes to the frequency domain.

    // In the frequency domain we are working with 1024 magnitudes and phases for every channel (left, right), if the fft size is 2048.
    while (frequencyDomain->timeDomainToFrequencyDomain(magnitudeLeft, magnitudeRight, phaseLeft,
                                                        phaseRight)) {
        // You can work with frequency domain data from this point.


        int start = 836;
        float maxVal = 0.0;
        int maxIndex = 0;
        for (int k = start; k < start+numberOfChannels; k ++) {
            if (magnitudeLeft[k] > treashold) {
                vals[k - start] = 1;
                //__android_log_print(ANDROID_LOG_DEBUG,"FREQ", "%d, %f", k*samplerate/frequencyDomain->fftSize, magnitudeLeft[k]  );

            } else {
                vals[k - start] = 0;
                //__android_log_print(ANDROID_LOG_DEBUG,"FREQ", "%d, %f, %s", k*samplerate/frequencyDomain->fftSize,magnitudeLeft[k], "MISS"  );

            }

        }
        memset(magnitudeLeft, 0, frequencyDomain->fftSize * sizeof(float));
        memset(magnitudeRight, 0, frequencyDomain->fftSize * sizeof(float));
        frequencyDomain->advance();
    }
    if (fifoOutputLastSample - fifoOutputFirstSample >= numberOfSamples) {
        SuperpoweredFloatToShortInt(fifoOutput + fifoOutputFirstSample * 2, audioInputOutput, (unsigned int)numberOfSamples);
        fifoOutputFirstSample += numberOfSamples;

        return true;
    }
    for (int i = 0; i < 3; i++) {
        for(int j = 0;j < 8; j++) {
            rgb[i] +=vals[j + i*8] * pow(2,j);
        }
    }
    jmethodID metodId = getEnv()->GetMethodID(mClass,"mColor","(III)V");
    if (metodId != NULL){
        getEnv()->CallVoidMethod(gClassLoader, metodId, rgb[0], rgb[1], rgb[2]);
    }
    return false;
}

extern "C" JNIEXPORT void
Java_se_livetsord_youth_MainActivity_FrequencyDomain(JNIEnv *__unused javaEnvironment,
                                                     jobject __unused obj, jint samplerate,
                                                     jint buffersize) {
    javaEnvironment->GetJavaVM(&gJvm);
    gClassLoader = javaEnvironment->NewGlobalRef(obj);
    frequencyDomain = new SuperpoweredFrequencyDomain(FFT_LOG_SIZE); // This will do the main "magic".
    stepSize = frequencyDomain->fftSize / 4; // The default overlap ratio is 4:1, so we will receive this amount of samples from the frequency domain in one step.

    // Frequency domain data goes into these buffers:
    magnitudeLeft = (float *) malloc(frequencyDomain->fftSize * sizeof(float));
    magnitudeRight = (float *) malloc(frequencyDomain->fftSize * sizeof(float));
    phaseLeft = (float *) malloc(frequencyDomain->fftSize * sizeof(float));
    phaseRight = (float *) malloc(frequencyDomain->fftSize * sizeof(float));
    // Time domain result goes into a FIFO (first-in, first-out) buffer
    fifoOutputFirstSample = fifoOutputLastSample = 0;
    fifoCapacity = stepSize * 100; // Let's make the fifo's size 100 times more than the step size, so we save memory bandwidth.
    fifoOutput = (float *) malloc(fifoCapacity * sizeof(float) * 2 + 128);

    inputBufferFloat = (float *) malloc(buffersize * sizeof(float) * 2 + 128);

    SuperpoweredCPU::setSustainedPerformanceMode(true);
    new SuperpoweredAndroidAudioIO(samplerate, buffersize, true, true, audioProcessing, NULL, -1,
                                   SL_ANDROID_STREAM_MEDIA,
                                   buffersize * 2); // Start audio input/output.
}

