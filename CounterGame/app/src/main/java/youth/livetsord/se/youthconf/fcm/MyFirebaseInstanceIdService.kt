package youth.livetsord.se.youthconf.fcm

/**
 * Created by Farez on 2017-07-13.
 *
 */


import android.content.Context
import android.os.AsyncTask
import android.util.Log
import com.github.kittinunf.fuel.Fuel

import com.google.firebase.iid.FirebaseInstanceId
import com.google.firebase.iid.FirebaseInstanceIdService
import com.google.firebase.messaging.FirebaseMessaging
import youth.livetsord.se.youthconf.Constants


class MyFirebaseInstanceIDService : FirebaseInstanceIdService() {

    /**
     * Called if InstanceID token is updated. This may occur if the security of
     * the previous token had been compromised. Note that this is called when the InstanceID token
     * is initially generated so this is where you would retrieve the token.
     */
    // [START refresh_token]
    override fun onTokenRefresh() {
        // Get updated InstanceID token.
        val refreshedToken = FirebaseInstanceId.getInstance().token
        Log.d(TAG, "Refreshed token: " + refreshedToken!!)

        // If you want to send messages to this application instance or
        // manage this apps subscriptions on the server side, send the
        // Instance ID token to your app server.
        if (getSharedPreferences(Constants.fcm.fcmSharedpref, Context.MODE_PRIVATE).getBoolean(Constants.fcm.isRegistered, false)) {
            val oldToken = getSharedPreferences(Constants.fcm.fcmSharedpref, Context.MODE_PRIVATE).getString(Constants.fcm.fcmToken, "")
            getSharedPreferences(Constants.fcm.fcmSharedpref, Context.MODE_PRIVATE).edit().putString(Constants.fcm.fcmToken, refreshedToken).apply()
            sendUpdatedRegistrationToServer(refreshedToken, oldToken)
        } else {
            sendRegistrationToServer(refreshedToken, "")
        }
    }
    // [END refresh_token]

    /**
     * Persist token to third-party servers.

     * Modify this method to associate the user's FCM InstanceID token with any server-side account
     * maintained by your application.

     * @param token The new token.
     */
    private fun sendRegistrationToServer(token: String, payload: String) {
        sendToServerAsynctask(token, payload, "user").execute()
        FirebaseMessaging.getInstance().subscribeToTopic("all")


    }

    private fun sendUpdatedRegistrationToServer(token: String, payload: String) {
        sendToServerAsynctask(token, payload, "user").execute()

    }

    companion object {

        private val TAG = "youthconf"
    }
}


class sendToServerAsynctask(val firebaseId: String, val payload: String, val type: String) : AsyncTask<String, Void, String>() {

    val host = "http://Youthappredirect.livetsord.se:8080/"


    override fun doInBackground(vararg p0: String?): String {
        var mResponse = ""
        Fuel.post(host + type + "/" + firebaseId).body(payload).response { request, response, result ->
            Log.d("Response", response.toString())
            mResponse = response.toString()
        }
        Log.d("Response", mResponse)
        return mResponse
    }

}
