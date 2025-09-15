package constants

import (
	"fmt"
	"os"
	"time"
)

const (
	FileSizeThreshold    = 1_000_000 // 1MB
	PNGType              = "image/png"
	JPGType              = "image/jpeg"
	PDFType              = "application/pdf"
	UnverifiedStatus     = "unverified"
	VerifiedStatus       = "verified"
	CTX_AUTH_PAYLOAD_KEY = "auth_payload"
	CTX_AUTH_TOKEN_KEY   = "auth_token"
	DefaultAvatar        = "https://res.cloudinary.com/dk8rlicon/image/upload/v1753881263/default-avatar-icon-of-social-media-user-vector_j8obqd.jpg"
	AvatarFolder         = "Avatars/"
	ResumeFolder         = "Resumes/"
	MaxCourseCategories  = 5
	ExpiredInterval      = 15 * time.Minute // 15 minute for development and testing
	NoRating             = 0
	MaxLimit             = 25
)

const (
	ReservedStatus       = "reserved"
	PendingPaymentStatus = "pending payment"
	ScheduledStatus      = "scheduled"
	CompletedStatus      = "completed"
	CancelledStatus      = "cancelled"
)

const (
	AdminRole   = iota + 1 // 1
	MentorRole             // 2
	StudentRole            // 3
)

const (
	RequestPerSecond = 10
	BurstSize        = 20
)

const (
	DefaultLimit  = 15
	DefaultPage   = 1
	DefaultLastID = 15
)

const (
	ForVerification = iota + 1 // 1
	ForReset
	ForLogin
)

const (
	StudentResource = iota + 1 // 1
	MentorResource
	CourseCategoryResource
	CourseResource
	CourseRequestResource
	ChatroomResource
	PaymentDetailResource
	PaymentMethodResource
	CourseRatingResource
	DiscountResource
	AdditionalCostResource
	AdminResource
)

const (
	CreatePermission = iota + 1 // 1
	ReadOwnPermission
	ReadAllPermission
	UpdateOwnPermission
	UpdateAllPermission
	DeleteOwnPermission
	DeleteAllPermission
)

func SendMentorAccEmailBody(email string, password string) string {
	loginURL := fmt.Sprintf("%s/mentor/login", os.Getenv("CLIENT_DOMAIN"))
	return fmt.Sprintf(`
    <!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your Account Credentials</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color: #f4f4f4;">
        <tr>
            <td style="padding: 20px 0;">
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="margin: 0 auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
                    
                    <!-- Header -->
                    <tr>
                        <td style="background-color: #2c3e50; padding: 30px; text-align: center; border-radius: 8px 8px 0 0;">
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: bold;">Privat-Unmei</h1>
                        </td>
                    </tr>
                    
                    <!-- Security Warning -->
                    <tr>
                        <td style="padding: 0 30px;">
                            <div style="margin: 20px 0 0 0; padding: 20px; background-color: #f8d7da; border: 2px solid #f5c6cb; border-radius: 5px;">
                                <h3 style="color: #721c24; margin: 0 0 10px 0; font-size: 16px;">🔒 SECURITY NOTICE</h3>
                                <p style="color: #721c24; font-size: 14px; margin: 0; line-height: 20px; font-weight: bold;">
                                    This email contains sensitive login credentials. Please read all security instructions below and change your password immediately after first login.
                                </p>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Main Content -->
                    <tr>
                        <td style="padding: 20px 30px 40px 30px;">
                            <h2 style="color: #2c3e50; margin: 0 0 20px 0; font-size: 24px;">Your Account Has Been Created</h2>
                            
                            <p style="color: #555555; font-size: 16px; line-height: 24px; margin: 0 0 30px 0;">
                                Your account has been successfully created. Below are your login credentials to access your Privat-Unmei account.
                            </p>
                            
                            <!-- Credentials Box -->
                            <div style="margin: 30px 0; padding: 25px; background-color: #f8f9fa; border: 2px solid #dee2e6; border-radius: 8px;">
                                <h3 style="color: #2c3e50; margin: 0 0 15px 0; font-size: 18px;">Your Login Credentials:</h3>
                                
                                <table style="width: 100%%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 10px 0; font-weight: bold; color: #2c3e50; width: 120px;">Email:</td>
                                        <td style="padding: 10px 15px; background-color: #ffffff; border: 1px solid #dee2e6; border-radius: 4px; font-family: 'Courier New', monospace; font-size: 16px; color: #2c3e50;">%s</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 10px 0; font-weight: bold; color: #2c3e50;">Password:</td>
                                        <td style="padding: 10px 15px; background-color: #ffffff; border: 1px solid #dee2e6; border-radius: 4px; font-family: 'Courier New', monospace; font-size: 16px; color: #2c3e50;">%s</td>
                                    </tr>
                                </table>
                            </div>
                            
                            <!-- Login Button -->
                            <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="%s" style="background-color: #27ae60; color: #ffffff; padding: 15px 30px; text-decoration: none; border-radius: 5px; font-size: 18px; font-weight: bold; display: inline-block; min-width: 200px;">
                                            Login to Your Account
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="color: #555555; font-size: 14px; line-height: 20px; margin: 20px 0; padding: 15px; background-color: #f8f9fa; border-left: 4px solid #27ae60; border-radius: 4px;">
                                <strong>Direct Login Link:</strong><br>
                                <a href="%s" style="color: #27ae60; word-break: break-all;">%s</a>
                            </p>
                            
                            <!-- Critical Security Instructions -->
                            <div style="margin: 30px 0; padding: 20px; background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 5px;">
                                <h3 style="color: #856404; margin: 0 0 15px 0; font-size: 16px;">🚨 IMMEDIATE ACTION REQUIRED:</h3>
                                <ol style="color: #856404; font-size: 14px; margin: 0; padding-left: 20px; line-height: 22px;">
                                    <li><strong>Login immediately</strong> and change your password</li>
                                    <li><strong>Delete this email</strong> after copying your credentials</li>
                                    <li><strong>Do not share</strong> these credentials with anyone</li>
                                    <li><strong>Use a strong, unique password</strong> when changing it</li>
                                </ol>
                            </div>
                        </td>
                    </tr>
                    
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
    `, email, password, loginURL, loginURL, loginURL)
}

func VerificationEmailBody(id string) string {
	verificationUrl := fmt.Sprintf("%s/verify/%s", os.Getenv("CLIENT_DOMAIN"), id)
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email Address</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color: #f4f4f4;">
        <tr>
            <td style="padding: 20px 0;">
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="margin: 0 auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
                    
                    <!-- Header -->
                    <tr>
                        <td style="background-color: #2c3e50; padding: 30px; text-align: center; border-radius: 8px 8px 0 0;">
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: bold;">Privat Unmei</h1>
                        </td>
                    </tr>
                    
                    <!-- Main Content -->
                    <tr>
                        <td style="padding: 40px 30px;">
                            <h2 style="color: #2c3e50; margin: 0 0 20px 0; font-size: 24px;">Welcome! Please verify your email</h2>
                            
                            <p style="color: #555555; font-size: 16px; line-height: 24px; margin: 0 0 20px 0;">
                                Thank you for creating an account with us! To complete your registration and secure your account, please verify your email address by clicking the button below.
                            </p>
                            
                            <!-- Verification Button -->
                            <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="%s" style="background-color: #3498db; color: #ffffff; padding: 15px 30px; text-decoration: none; border-radius: 5px; font-size: 18px; font-weight: bold; display: inline-block; min-width: 200px;">
                                            Verify Email Address
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="color: #555555; font-size: 14px; line-height: 20px; margin: 20px 0; padding: 15px; background-color: #f8f9fa; border-left: 4px solid #3498db; border-radius: 4px;">
                                <strong>Can't click the button?</strong><br>
                                Copy and paste this link into your browser:<br>
                                <a href="%s" style="color: #3498db; word-break: break-all;">%s</a>
                            </p>
                            
                            <div style="margin: 30px 0; padding: 20px; background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 5px;">
                                <h3 style="color: #856404; margin: 0 0 10px 0; font-size: 16px;">Important Information:</h3>
                                <ul style="color: #856404; font-size: 14px; margin: 0; padding-left: 20px;">
                                    <li>This verification link will expire in <strong>72 hours</strong></li>
                                    <li>If you didn't create an account with us, please ignore this email</li>
                                    <li>For security reasons, do not share this verification link with anyone</li>
                                </ul>
                            </div>
	`, verificationUrl, verificationUrl, verificationUrl)
}

func ResetEmailBody(token string) string {
	resetURL := fmt.Sprintf("%s/reset/%s", os.Getenv("CLIENT_DOMAIN"), token)
	return fmt.Sprintf(`
    <!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color: #f4f4f4;">
        <tr>
            <td style="padding: 20px 0;">
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="margin: 0 auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
                    
                    <!-- Header -->
                    <tr>
                        <td style="background-color: #2c3e50; padding: 30px; text-align: center; border-radius: 8px 8px 0 0;">
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: bold;">Privat Unmei</h1>
                        </td>
                    </tr>
                    
                    <!-- Main Content -->
                    <tr>
                        <td style="padding: 40px 30px;">
                            <h2 style="color: #2c3e50; margin: 0 0 20px 0; font-size: 24px;">Reset Your Password</h2>
                            
                            <p style="color: #555555; font-size: 16px; line-height: 24px; margin: 0 0 20px 0;">
                                We received a request to reset your password for your Privat Unmei account. Click the button below to create a new password.
                            </p>
                            
                            <!-- Reset Password Button -->
                            <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="%s" style="background-color: #e74c3c; color: #ffffff; padding: 15px 30px; text-decoration: none; border-radius: 5px; font-size: 18px; font-weight: bold; display: inline-block; min-width: 200px;">
                                            Reset Password
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="color: #555555; font-size: 14px; line-height: 20px; margin: 20px 0; padding: 15px; background-color: #f8f9fa; border-left: 4px solid #e74c3c; border-radius: 4px;">
                                <strong>Can't click the button?</strong><br>
                                Copy and paste this link into your browser:<br>
                                <a href="%s" style="color: #e74c3c; word-break: break-all;">%s</a>
                            </p>
                            
                            <div style="margin: 30px 0; padding: 20px; background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 5px;">
                                <h3 style="color: #856404; margin: 0 0 10px 0; font-size: 16px;">Important Security Information:</h3>
                                <ul style="color: #856404; font-size: 14px; margin: 0; padding-left: 20px;">
                                    <li>This password reset link will expire in <strong>72 hours</strong></li>
                                    <li>If you didn't request a password reset, please ignore this email</li>
                                    <li>Your password will remain unchanged if you don't click the link</li>
                                    <li>For security reasons, do not share this reset link with anyone</li>
                                </ul>
                            </div>
                        </td>
                    </tr>
                    
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
    `, resetURL, resetURL, resetURL)
}
