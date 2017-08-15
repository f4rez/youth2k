package youth.livetsord.se.countergame

import android.content.Context
import android.os.AsyncTask
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.widget.RadioButton
import android.widget.Toast
import com.github.kittinunf.fuel.Fuel
import com.google.firebase.iid.FirebaseInstanceId


import com.google.gson.GsonBuilder

import kotlinx.android.synthetic.main.activity_register.*
import youth.livetsord.se.countergame.fcm.Userpost


class Register : AppCompatActivity(), RestCallback {





    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register)
        Log.d("Register", "hej")
        save_selection.setOnClickListener {
            Log.d("Register", "save_selection clicked")
            if (radio_buttons.checkedRadioButtonId >= 0) {
                val rb = findViewById(radio_buttons.checkedRadioButtonId) as RadioButton
                if (name_edittext.length() in 5..19) {
                    val team = rb.text.toString()
                    val name = name_edittext.text.toString()
                    val up = Userpost(name, team)
                    val builder = GsonBuilder()
                    val gson = builder.create()
                    val payload = gson.toJson(up)
                    AsyncRestSender(this, FirebaseInstanceId.getInstance().token.toString(), payload, "user").execute()

                } else {
                    Toast.makeText(applicationContext, "Du måste ha ett namn som är mellan 5-20 tecken", Toast.LENGTH_SHORT).show()
                }
            } else {
                Toast.makeText(applicationContext, "Välj ett lag", Toast.LENGTH_SHORT).show()
            }


        }

    }

    override fun callback(code: Int) {
        getSharedPreferences(Constants.fcm.fcmSharedpref, Context.MODE_PRIVATE).edit().putBoolean(Constants.fcm.isRegistered, true).apply()

            when (code) {
            Constants.codes.success -> finish()
            Constants.codes.unsuccessful -> {}
            Constants.codes.nameAllreadyInUse -> Toast.makeText(applicationContext, "Namnet är redan upptaget", Toast.LENGTH_SHORT).show()
        }

    }

}


class AsyncRestSender(val c: youth.livetsord.se.countergame.RestCallback, val firebaseId : String, val payload : String, val type : String) : AsyncTask<String, Void, String>() {



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

