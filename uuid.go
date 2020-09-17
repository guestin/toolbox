package mob

import (
	"github.com/google/uuid"
	"strings"
)

//创建随机uuid id
func GenRandomUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
