package youth.livetsord.se.youthconf

import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.graphics.Bitmap
import android.net.ConnectivityManager
import android.os.Build
import android.os.Bundle
import android.os.Environment
import android.support.annotation.RequiresApi
import android.support.v4.app.ActivityCompat
import android.support.v4.widget.SwipeRefreshLayout
import android.support.v7.app.AppCompatActivity
import android.util.Log
import android.view.View
import android.view.animation.AnimationUtils
import android.webkit.WebResourceRequest
import android.webkit.WebSettings
import android.webkit.WebSettings.LOAD_CACHE_ELSE_NETWORK
import android.webkit.WebView
import android.webkit.WebViewClient
import android.widget.Toast
import com.squareup.picasso.Picasso
import kotlinx.android.synthetic.main.start_page.*
import org.jetbrains.anko.doAsync
import org.jetbrains.anko.uiThread
import java.io.File
import java.io.FileOutputStream


class MainActivity : AppCompatActivity(), SwipeRefreshLayout.OnRefreshListener {

    var count = 0
    val tag = "youthconf"
    var loading = false

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContentView(R.layout.start_page)
        Log.d(tag, "onCreate begin")
        var webSettings = mWebView.settings
        webSettings.javaScriptEnabled = true
        webSettings.setAppCachePath(applicationContext.cacheDir.absolutePath)
        webSettings.allowFileAccess = true
        webSettings.setAppCacheEnabled(true)
        webSettings.cacheMode = WebSettings.LOAD_CACHE_ELSE_NETWORK
        webSettings = mWebView2.settings
        webSettings.javaScriptEnabled = true
        webSettings.setAppCachePath(applicationContext.cacheDir.absolutePath)
        webSettings.allowFileAccess = true
        webSettings.setAppCacheEnabled(true)
        webSettings.cacheMode = WebSettings.LOAD_CACHE_ELSE_NETWORK
        // loading offline
        if (!isNetworkAvailable()) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK



        home.setOnClickListener {
            if (!loading) {
                if (!isNetworkAvailable()) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK
                loadUrl(Constants.Urls.home)
                // loading offline
                unHighlight()
                home.setImageResource(R.drawable.hem_vit)
            }
        }

        speakers.setOnClickListener {
            if (!loading) {
                if (!isNetworkAvailable()) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK
                loadUrl(Constants.Urls.speakers)
                unHighlight()
                speakers.setImageResource(R.drawable.talare_vit)
            }
        }

        schedule.setOnClickListener {
            if (!loading) {
                if (!isNetworkAvailable()) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK
                loadUrl(Constants.Urls.schedule)
                unHighlight()
                schedule.setImageResource(R.drawable.schema_vit)
            }
        }

        live.setOnClickListener {
            if (!loading) {
                if (!isNetworkAvailable()) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK
                loadUrl(Constants.Urls.downloads)
                unHighlight()
                live.setImageResource(R.drawable.dwnlds_vit)
            }
        }
        val mWebViewClient = (object : WebViewClient() {

            override fun onPageFinished(view: WebView, url: String) {
                if (url.endsWith(".png") || url.endsWith(".jpg")) {
                    Log.d(tag, "Should download image at: " + url)

                } else if (swipeRefresh1.isRefreshing || swipeRefresh2.isRefreshing) {
                    swipeRefresh1.isRefreshing = false
                    swipeRefresh2.isRefreshing = false

                } else if (splashscreen.visibility == View.VISIBLE) {
                    val animFadeOut = AnimationUtils.loadAnimation(applicationContext, android.R.anim.fade_out)
                    val animFadeIn = AnimationUtils.loadAnimation(applicationContext, android.R.anim.fade_in)
                    fadeBottomBar(url)
                    progressBar.animation = animFadeOut
                    splashscreen.animation = animFadeOut
                    container.animation = animFadeIn
                    container.visibility = View.VISIBLE
                    splashscreen.visibility = View.GONE
                    progressBar.visibility = View.GONE
                    Log.d(tag, "Incrementing count from splashscreen")
                    count++
                } else {
                    switchWebView(url)
                    Log.d(tag, "Incrementing count from default")
                    count++
                }

                loading = false
                stopSpinner()
            }

            override fun shouldOverrideUrlLoading(view: WebView?, url: String?): Boolean {
                if (url!!.endsWith(suffix = ".png", ignoreCase = true) || url.endsWith(suffix = ".jpg", ignoreCase = true)) {
                    if (isStoragePermissionGranted()) {
                        Log.v(tag, "Permission is granted")
                        doAsync {
                            saveImage(url)
                            uiThread {
                                Toast.makeText(applicationContext, "Laddade ner bilden till ditt galleri", Toast.LENGTH_SHORT).show()
                            }
                        }
                    }
                    return true
                }
                return false
            }

            @RequiresApi(Build.VERSION_CODES.LOLLIPOP)
            override fun shouldOverrideUrlLoading(view: WebView?, request: WebResourceRequest?): Boolean {

                val url = request!!.url.toString()

                if (url.endsWith(".png", true) || url.endsWith(".jpg", true)) {
                    if (isStoragePermissionGranted()) {
                        Log.v(tag, "Permission is granted")
                        doAsync {
                            saveImage(url)
                            uiThread {
                                Toast.makeText(applicationContext, "Laddade ner bilden till ditt galleri", Toast.LENGTH_SHORT).show()
                            }
                        }
                    }

                    return true
                }
                return false
            }
        })
        mWebView.setWebViewClient(mWebViewClient)
        mWebView2.setWebViewClient(mWebViewClient)

        swipeRefresh1.setOnRefreshListener(this)
        swipeRefresh2.setOnRefreshListener(this)

        onNewIntent(intent)


    }


    override fun onNewIntent(intent: Intent?) {
        Log.d(tag, "New intent arrived")
        if (intent != null) {
            if (intent.extras != null && intent.extras.get("link") != null) {
                Log.d(tag, intent.extras.get("link") as String)
                loadUrl(intent.extras.get("link") as String)
                unHighlight()
            } else {
                loadUrl(Constants.Urls.home)

            }
            setIntent(intent)
        }
    }


    private fun unHighlight() {
        home.setImageResource(R.drawable.hem)
        speakers.setImageResource(R.drawable.talare)
        schedule.setImageResource(R.drawable.schema)
        live.setImageResource(R.drawable.dwnlds)
    }


    private fun switchWebView(url: String) {
        val animFadeOut = AnimationUtils.loadAnimation(applicationContext, R.anim.fade_out)
        val animFadeIn = AnimationUtils.loadAnimation(applicationContext, R.anim.fade_in)
        if (count % 2 == 1 && mWebView2.url != null) {
            Log.d(tag, "Switching from mWebView to mWebView2, url to load = $url, count = $count")
            mWebView.animation = animFadeOut
            mWebView2.animation = animFadeIn
            mWebView2.visibility = View.VISIBLE
            mWebView.visibility = View.GONE
            swipeRefresh1.visibility = View.GONE
            swipeRefresh2.visibility = View.VISIBLE

        } else if (mWebView.url != null) {
            Log.d(tag, "Switching from mWebView2 to mWebView, url to load = $url, count = $count")
            mWebView2.animation = animFadeOut
            mWebView.animation = animFadeIn
            mWebView.visibility = View.VISIBLE
            mWebView2.visibility = View.GONE
            swipeRefresh1.visibility = View.VISIBLE
            swipeRefresh2.visibility = View.GONE
        }
        loading = false

        Log.d(tag, " " + Constants.Urls.home + " " + (fade.visibility == View.VISIBLE))
        fadeBottomBar(url)
        Log.d(tag, "Web1: " + mWebView.url + ", Web2: " + mWebView2.url)


    }

    private fun fadeBottomBar(url: String) {
        val animFadeOut = AnimationUtils.loadAnimation(applicationContext, R.anim.fade_out)
        val animFadeIn = AnimationUtils.loadAnimation(applicationContext, R.anim.fade_in)
        if (url == Constants.Urls.home) {
            fade.animation = animFadeOut
            bottombarColor.animation = animFadeOut
            fade.visibility = View.GONE
            bottombarColor.visibility = View.GONE
        } else if (fade.visibility == View.GONE) {
            fade.animation = animFadeIn
            bottombarColor.animation = animFadeIn
            fade.visibility = View.VISIBLE
            bottombarColor.visibility = View.VISIBLE
        }
    }

    private fun startSpinner() {
        if (progressBar.visibility != View.VISIBLE) {
            val animFadeIn = AnimationUtils.loadAnimation(applicationContext, R.anim.fade_in)
            Log.d(tag, "Starting spinner")
            progressBarLoadingWeb.animation = animFadeIn
            progressBarLoadingWeb.visibility = View.VISIBLE
        }

    }

    private fun stopSpinner() {
        val animFadeOut = AnimationUtils.loadAnimation(applicationContext, R.anim.fade_out)
        Log.d(tag, "Stop spinner")

        progressBarLoadingWeb.animation = animFadeOut
        progressBarLoadingWeb.visibility = View.GONE

    }

    private fun loadUrl(url: String) {
        Log.d(tag, "Should load url: " + url)
        loading = true
        startSpinner()

        if (count % 2 == 0) {
            Log.d(tag, "loadUrl - current: " + mWebView2.url + ", load: " + url)
            if (url != mWebView2.url) {
                mWebView.loadUrl(url)
            } else {
                loading = false
                stopSpinner()
            }
        } else {
            Log.d(tag, "loadUrl - current2: " + mWebView.url + ", load: " + url)

            if (url != mWebView.url) {
                mWebView2.loadUrl(url)
            } else {
                loading = false
                stopSpinner()
            }

        }


    }


    private fun isNetworkAvailable(): Boolean {
        val connectivityManager = getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager
        val activeNetworkInfo = connectivityManager.activeNetworkInfo
        return activeNetworkInfo != null && activeNetworkInfo.isConnected
    }

    override fun onRefresh() {
        if (count % 2 == 0) {
            mWebView2.reload()
        } else {
            mWebView.reload()
        }
    }

    private fun saveImage(url: String) {
        val bitmap = Picasso.with(this).load(url).get()
        val fileName = url.substring(url.dropLast(1).lastIndexOf("/") + 1, url.length)


        val direct = File(Environment.getExternalStoragePublicDirectory("/Pictures/Youth").toURI())

        if (!direct.exists()) {
            Log.d(tag, "Creating folder Picturs/Youth")
            val wallpaperDirectory = File(Environment.getExternalStoragePublicDirectory("/Pictures/Youth").toURI())
            Log.d(tag, "  " + wallpaperDirectory.mkdirs())
        }

        val file = File(File(Environment.getExternalStoragePublicDirectory("/Pictures/Youth").toURI()), fileName)
        if (file.exists()) {
            file.delete()
        }
        try {
            val out = FileOutputStream(file)
            bitmap.compress(Bitmap.CompressFormat.PNG, 100, out)
            out.flush()
            out.close()
            val intent = Intent(android.content.Intent.ACTION_VIEW)
            intent.setDataAndType(android.support.v4.content.FileProvider.getUriForFile(this, packageName + ".provider", file), "image/*")
            intent.addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION)
            Log.d(tag, " " + intent.type)

            val pendingIntent = android.app.PendingIntent.getActivity(this, Constants.recuest.openLink, intent,
                    android.app.PendingIntent.FLAG_CANCEL_CURRENT)

            val notificationBuilder = android.support.v7.app.NotificationCompat.Builder(this)
                    .setSmallIcon(youth.livetsord.se.youthconf.R.drawable.push_notification_icon)
                    .setContentTitle("Youth")
                    .setContentText("Här är din bild!")
                    .setAutoCancel(true)
                    .setContentIntent(pendingIntent)
            val notificationManager = getSystemService(android.content.Context.NOTIFICATION_SERVICE) as android.app.NotificationManager

            notificationManager.notify(0 /* ID of notification */, notificationBuilder.build())
        } catch (e: Exception) {
            Log.e(tag, "Error saving image: " + e.printStackTrace())
        }

    }

    fun isStoragePermissionGranted(): Boolean {
        return if (Build.VERSION.SDK_INT >= 23) {
            if (checkSelfPermission(android.Manifest.permission.WRITE_EXTERNAL_STORAGE) == PackageManager.PERMISSION_GRANTED) {
                Log.v(tag, "Permission is granted")
                true
            } else {
                Log.v(tag, "Permission is revoked")
                ActivityCompat.requestPermissions(this, arrayOf(android.Manifest.permission.WRITE_EXTERNAL_STORAGE), 1)
                false
            }
        } else { //permission is automatically granted on sdk<23 upon installation
            Log.v(tag, "Permission is granted")
            true
        }
    }

    override fun onRequestPermissionsResult(requestCode: Int, permissions: Array<String>, grantResults: IntArray) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (grantResults[0] == PackageManager.PERMISSION_GRANTED) {
            Log.v(tag, "Permission: " + permissions[0] + "was " + grantResults[0])
        }
    }
}




