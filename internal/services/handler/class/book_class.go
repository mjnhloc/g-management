package class

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *HTTPHandler) BookClass(c *gin.Context) {
	userRole, _ := c.Get("user_role")
	if userRole != "member" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only members can book classes"})
		return
	}
	memberIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user ID"})
		return
	}
	var memberID int64
	switch v := memberIDRaw.(type) {
	case int64:
		memberID = v
	case int:
		memberID = int64(v)
	case float64:
		memberID = int64(v)
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		memberID = id
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}
	classID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema: h.graphql,
		RequestString: `
			mutation BookClass($class_id: Int!, $member_id: Int!) {
				book_class(class_id: $class_id, member_id: $member_id) {
					id
					name
					schedule
					description
					max_capacity
				}
			}
		`,
		VariableValues: map[string]interface{}{
			"class_id":  int(classID),
			"member_id": int(memberID),
		},
		Context: c.Request.Context(),
	})

	if len(result.Errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Errors[0].Message})
		return
	}

	c.JSON(http.StatusOK, result.Data)
}
