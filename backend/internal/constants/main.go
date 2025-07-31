package constants

import (
	"fmt"
	"os"
)

const (
	FileSizeThreshold = 1_000_000 // 1MB
	PNGType           = "image/png"
	JPGType           = "image/jpeg"
	PDFType           = "application/pdf"
	UnverifiedStatus  = "unverified"
	VerifiedStatus    = "verified"
	DefaultAvatar     = "https://res.cloudinary.com/dk8rlicon/image/upload/v1753881263/default-avatar-icon-of-social-media-user-vector_j8obqd.jpg"
)

const (
	StudentRole = 1
	AdminRole   = 2
	MentorRole  = 3
)

const (
	ForVerification = 1
	ForReset        = 2
	ForLogin        = 3
)

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
                                    <li>This verification link will expire in <strong>24 hours</strong></li>
                                    <li>If you didn't create an account with us, please ignore this email</li>
                                    <li>For security reasons, do not share this verification link with anyone</li>
                                </ul>
                            </div>
	`, verificationUrl, verificationUrl, verificationUrl)
}
