package main
import (
    "fmt"
    "os"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/kms"
)
func main() {
    keyID := "arn:aws:kms:us-east-1:123456789012:key/your-kms-key-id"
    sess := session.Must(session.NewSession())
    svc := kms.New(sess)
    // Encrypt
    plaintext := []byte("Top Secret Cloud Data")
    encryptInput := &kms.EncryptInput{
        KeyId:     &keyID,
        Plaintext: plaintext,
    }
    encryptOutput, err := svc.Encrypt(encryptInput)
    if err != nil {
        fmt.Println("Encrypt error:", err)
        os.Exit(1)
    }
    fmt.Println("Encrypted Data:", encryptOutput.CiphertextBlob)
    // Decrypt
    decryptInput := &kms.DecryptInput{
        CiphertextBlob: encryptOutput.CiphertextBlob,
    }
    decryptOutput, err := svc.Decrypt(decryptInput)
    if err != nil {
        fmt.Println("Decrypt error:", err)
        os.Exit(1)
    }
    fmt.Println("Decrypted Text:", string(decryptOutput.Plaintext))
}
