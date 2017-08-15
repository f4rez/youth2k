package youth.livetsord.se.countergame.fcm

import youth.livetsord.se.countergame.MainActivity


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

        // TODO(developer): Handle FCM messages here.
        // Not getting messages here? See why this may be: https://goo.gl/39bRNJ
        android.util.Log.d(youth.livetsord.se.countergame.fcm.MyFirebaseMessagingService.Companion.TAG, "From: " + remoteMessage!!.from)

        // Check if message contains a data payload.
        if (remoteMessage.data.isNotEmpty()) {
            android.util.Log.d(youth.livetsord.se.countergame.fcm.MyFirebaseMessagingService.Companion.TAG, "Message data payload: " + remoteMessage.data)

        }

        // Check if message contains a notification payload.
        if (remoteMessage.notification != null) {
            android.util.Log.d(youth.livetsord.se.countergame.fcm.MyFirebaseMessagingService.Companion.TAG, "Message Notification Body: " + remoteMessage.notification.body)
        }

        // Also if you intend on generating your own notifications as a result of a received FCM
        // message, here is where that should be initiated. See sendNotification method below.
    }
    // [END receive_message]


    /**
     * Handle time allotted to BroadcastReceivers.
     */
    private fun handleNow() {
        android.util.Log.d(youth.livetsord.se.countergame.fcm.MyFirebaseMessagingService.Companion.TAG, "Short lived task is done.")
    }

    /**
     * Create and show a simple notification containing the received FCM message.

     * @param messageBody FCM message body received.
     */
    private fun sendNotification(messageBody: String) {
        val intent = android.content.Intent(this, MainActivity::class.java)
        intent.addFlags(android.content.Intent.FLAG_ACTIVITY_CLEAR_TOP)
        val pendingIntent = android.app.PendingIntent.getActivity(this, 0 /* Request code */, intent,
                android.app.PendingIntent.FLAG_ONE_SHOT)

        val defaultSoundUri = android.media.RingtoneManager.getDefaultUri(android.media.RingtoneManager.TYPE_NOTIFICATION)
        val notificationBuilder = android.support.v7.app.NotificationCompat.Builder(this)
                .setSmallIcon(youth.livetsord.se.countergame.R.drawable.ic_arrow_downward_black_48dp)
                .setContentTitle("FCM Message")
                .setContentText(messageBody)
                .setAutoCancel(true)
                .setSound(defaultSoundUri)
                .setContentIntent(pendingIntent)

        val notificationManager = getSystemService(android.content.Context.NOTIFICATION_SERVICE) as android.app.NotificationManager

        notificationManager.notify(0 /* ID of notification */, notificationBuilder.build())
    }

    companion object {

        private val TAG = "MyFirebaseMsgService"
    }
}