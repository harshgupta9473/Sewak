package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"crypto/rand"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	otp := make([]byte, length)

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		otp[i] = charSet[index.Int64()]
	}

	return string(otp), nil
}



func SendOtp(mobile string, otp string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhone := os.Getenv("TWILIO_PHONE_NUMBER")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	params := &openapi.CreateMessageParams{}
	params.SetTo(mobile)
	params.SetFrom(twilioPhone)
	params.SetBody(fmt.Sprintf("Your OTP is: %s", otp))

	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		return fmt.Errorf("error sending OTP: %v", err)
	}

	log.Printf("OTP sent: SID %s", *resp.Sid)
	return nil

}