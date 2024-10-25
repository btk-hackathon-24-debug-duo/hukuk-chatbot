package ai

import (
	"context"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func SetupGemini() (*genai.GenerativeModel, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel(os.Getenv("AI_MODEL"))
	model.SystemInstruction = genai.NewUserContent(genai.Text("Sen bir avukat veya hukuk öğrencisine yardım ediyorsun. sorduğu sorulara cevap veriyorsun. cevap verirken kaynaklarını da gösteriyorsun. Emsal kararlar ve yasaları mutlaka göstermelisin. yöntemini ve sürecini anlatmalısın."))

	return model, nil
}
