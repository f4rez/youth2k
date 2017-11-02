package youth.livetsord.se.youthconf.fcm

import android.content.Intent
import youth.livetsord.se.youthconf.Constants
import youth.livetsord.se.youthconf.MainActivity


/**
 * Created by Farez on 2017-07-13.
 */

class MyFirebaseMessagingService : com.google.firebase.messaging.FirebaseMessagingService() {

    /**
     * Called when message is received.

     * @param remoteMessage Object representing the message received from Firebase Cloud Messaging.
     */
    // [START receive_message]
    override fun onMessageReceived(remoteMessage: com.google.firebase.messaging.RemoteMessage?) {
        // [START_EXCLUDE]
        // There are two types of messages data messages and notification messages. Data messages are handled
        // here in onMessageReceived whether the app is in the foreground or background. Data messages are the type
        // traditionally used with GCM. Notification messages are only received here in onMessageReceived when the app
        // is in the foreground. When the app is in the background an automatically generated notification is displayed.
        // When the user taps on the notification they are returned to the app. Messages containing both notification
        // and data payloads are treated as notification messages. The Firebase console always sends notification
        // messages. For more see: https://firebase.google.com/docs/cloud-messaging/concept-options
        // [END_EXCLUDE]


        android.util.Log.d(TAG, "From: " + remoteMessage!!.from)

        // Check if message contains a data payload.
        if (remoteMessage.data.isNotEmpty()) {
            android.util.Log.d(TAG, "Message data payload: " + remoteMessage.data)
            sendNotification(remoteMessage.notification.body, remoteMessage.data.getValue("link"))

        }

        // Check if message contains a notification payload.
        if (remoteMessage.notification != null) {
            android.util.Log.d(TAG, "Message Notification Body: " + remoteMessage.notification.body)
            sendNotification(remoteMessage.notification.body, remoteMessage.data.getValue("link"))
        }

        // Also if you intend on generating your own notifications as a result of a received FCM
        // message, here is where that should be initiated. See sendNotification method below.
    }
    // [END receive_message]


    /**
     * Handle time allotted to BroadcastReceivers.
     */
    private fun handleNow() {
        android.util.Log.d(TAG, "Short lived task is done.")
    }

    /**
     * Create and show a simple notification containing the received FCM message.

     * @param messageBody FCM message body received.
     */
    private fun sendNotification(messageBody: String, link: String?) {
        val intent = android.content.Intent(this, MainActivity::class.java)
        intent.addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP)
        intent.addFlags(Intent.FLAG_ACTIVITY_SINGLE_TOP)
        intent.putExtra("link", link)
        val pendingIntent = android.app.PendingIntent.getActivity(this, Constants.recuest.openLink, intent,
                android.app.PendingIntent.FLAG_CANCEL_CURRENT)


        val defaultSoundUri = android.media.RingtoneManager.getDefaultUri(android.media.RingtoneManager.TYPE_NOTIFICATION)
        val notificationBuilder = android.support.v7.app.NotificationCompat.Builder(this)
                .setSmallIcon(youth.livetsord.se.youthconf.R.drawable.push_notification_icon)
                .setContentTitle("Youth")
                .setContentText(messageBody)
                .setAutoCancel(true)
                .setSound(defaultSoundUri)
                .setContentIntent(pendingIntent)

        val notificationManager = getSystemService(android.content.Context.NOTIFICATION_SERVICE) as android.app.NotificationManager

        notificationManager.notify(0 /* ID of notification */, notificationBuilder.build())
    }

    companion object {

        private val TAG = "youthconf"
    }
}