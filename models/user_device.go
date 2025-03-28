package models

// UserDevice represents a user's device information for push notifications.
type UserDevice struct {
	DeviceID  string `json:"deviceId" bson:"device_id" binding:"required"`
	DeviceName string `json:"deviceName" bson:"device_name" binding:"required"`
	FcmToken  string `json:"fcmToken" bson:"fcm_token" binding:"required"`
	DeviceTimezone *string `json:"deviceTimezone" bson:"device_timezone"`
}