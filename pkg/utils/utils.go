package utils

import (
	"context"
	"fmt"

	"github.com/eiannone/keyboard"
)

// SetupEscListener sets up a goroutine to listen for ESC key presses to cancel streaming
func SetupEscListener(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		err := keyboard.Open()
		if err != nil {
			return
		}
		defer keyboard.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				char, key, err := keyboard.GetKey()
				if err != nil {
					continue
				}

				// Check for ESC key
				if key == keyboard.KeyEsc || char == 27 {
					fmt.Print("\n🛑 Stream stopped by user (ESC pressed)\n")
					cancel()
					return
				}
			}
		}
	}()
}