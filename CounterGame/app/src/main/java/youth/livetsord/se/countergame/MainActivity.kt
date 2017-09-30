package youth.livetsord.se.countergame

import android.content.Context
import android.content.Intent
import android.os.AsyncTask
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Toast
import com.github.kittinunf.fuel.Fuel
import kotlinx.android.synthetic.main.activity_main.*
import com.google.firebase.iid.FirebaseInstanceId
import com.google.gson.GsonBuilder
import kotlinx.android.synthetic.main.start_page.*
import android.webkit.WebSettings
import android.webkit.WebView
import android.webkit.WebViewClient
import android.view.animation.AnimationUtils
import android.net.ConnectivityManager
import android.webkit.WebSettings.LOAD_CACHE_ELSE_NETWORK


class MainActivity : AppCompatActivity() {


    var counter = 0
    var user = "Farez"
    var team = "green"


    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.start_page)
        val webSettings = mWebView.settings
        webSettings.javaScriptEnabled = true
        webSettings.setAppCachePath(applicationContext.cacheDir.absolutePath)
        webSettings.allowFileAccess = true
        webSettings.setAppCacheEnabled(true)
        webSettings.cacheMode =WebSettings.LOAD_DEFAULT
        // loading offline
        if ( !isNetworkAvailable() ) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK

        mWebView.loadUrl("http://youth.livetsord.se")

        speakers.setOnClickListener {
            if ( !isNetworkAvailable() ) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK

            mWebView.loadUrl("http://youth.livetsord.se")
            // loading offline

        }

        schedule.setOnClickListener {
            if ( !isNetworkAvailable() ) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK

            mWebView.loadUrl("http://livetsord.se")
            // loading offline
        }

        live.setOnClickListener {
            if ( !isNetworkAvailable() ) mWebView.settings.cacheMode = LOAD_CACHE_ELSE_NETWORK
            mWebView.loadUrl("http://google.se")
            // loading offline
        }

        mWebView.setWebViewClient(object : WebViewClient() {

            override fun onPageFinished(view: WebView, url: String) {
                val animFadeOut = AnimationUtils.loadAnimation(applicationContext, android.R.anim.fade_out)
                val animFadeIn = AnimationUtils.loadAnimation(applicationContext, android.R.anim.fade_in)
                splashscreen.animation = animFadeOut
                container.animation = animFadeIn
                container.visibility = View.VISIBLE
                splashscreen.visibility = View.GONE
            }
        })




        /*if (!getSharedPreferences(Constants.fcm.fcmSharedpref, Context.MODE_PRIVATE).getBoolean(Constants.fcm.isRegistered, false)) {
            val intent = Intent(this, Register::class.java)
            startActivity(intent)
        }

        team1_upvote.setOnClickListener(this)
        team2_upvote.setOnClickListener(this)
        team3_upvote.setOnClickListener(this)
        team1_downvote.setOnClickListener(this)
        team2_downvote.setOnClickListener(this)
        team3_downvote.setOnClickListener(this)

*/


    }

    private fun isNetworkAvailable(): Boolean {
        val connectivityManager = getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager
        val activeNetworkInfo = connectivityManager.activeNetworkInfo
        return activeNetworkInfo != null && activeNetworkInfo.isConnected
    }

/*
    override fun onClick(p0: View?) {
        val l = p0!!.tag.toString().split(" ")
        var team = ""
        when (l[0]) {
            "team1" -> team = getString(R.id.team1)
            "team2" -> team = getString(R.id.team2)
            "team3" -> team = getString(R.id.team3)
        }
        val vote = Vote(team, l[1].toInt())
        val builder = GsonBuilder()
        val gson = builder.create()
        val payload = gson.toJson(vote)
        SendCountToServer(this,FirebaseInstanceId.getInstance().id,payload,"vote").execute()



    }

    override fun callback(code: Int) {

    }


    inner class Vote(val team: String, val count: Int) {}

    inner class VoteResponse(val team1: String, val team1_count: Int, val team2: String, val team2_count: Int, val team3: String, val team3_count: Int) {}

    inner class SendCountToServer(val c: youth.livetsord.se.countergame.RestCallback, val firebaseId : String, val payload : String, val type : String) : AsyncTask<String, Void, String>() {

        override fun doInBackground(vararg p0: String?): String {
            var mResponse = ""
            Fuel.put(Constants.host.host + type + "/" + firebaseId).body(payload).response { request, response, result ->
                Log.d("tja", response.toString() + "\n" + request.toString())
                mResponse = response.toString()
            }
            Log.d("tja", mResponse)
            return mResponse
        }

        override fun onPostExecute(result: String?) {
            super.onPostExecute(result)
            c.callback(1)
        }

    }

*/

}


