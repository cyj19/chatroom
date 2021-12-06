/**
 * @Author: cyj19
 * @Date: 2021/12/6 14:18
 */

// 敏感词过滤

package logic

import (
	"github.com/cyj19/chatroom/global"
	"strings"
)

func FilterSensitiveWords(content string) string {
	for _, word := range global.SensitiveWords {
		content = strings.ReplaceAll(content, word, "**")
	}
	return content
}
