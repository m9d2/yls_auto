package model

type User struct {
	AvatarUrl       string   `json:"avatar_url"`
	NickName        string   `json:"nick_name"`
	PhoneNumber     string   `json:"phone_number"`
	CurrentScoreSum int      `json:"currentScoreSum"`
	Grade           string   `json:"grade"`
	GradeCode       string   `json:"gradeCode"`
	ConsumerID      string   `json:"consumerID"`
	IsMember        int      `json:"is_member"`
	Text            string   `json:"text"`
	TextArr         []string `json:"text_arr"`
	QrLink          string   `json:"qr_link"`
	QrLinkArr       struct {
		White  string `json:"white"`
		Black  string `json:"black"`
		Golden string `json:"golden"`
	} `json:"qr_link_arr"`
	MemberTips            string        `json:"memberTips"`
	ProgressBar           float64       `json:"progressBar"`
	ExpireDateText        string        `json:"expireDateText"`
	WelcomeMessage        string        `json:"welcome_message"`
	MembershipBenefitsTip []string      `json:"membership_benefits_tip"`
	MessageTip            []interface{} `json:"message_tip"`
	IsBirthdayMonth       bool          `json:"is_birthday_month"`
	Birthday              string        `json:"birthday"`
	ExpireScore           int           `json:"expireScore"`
	ExpireScoreDate       string        `json:"expireScoreDate"`
	ExpirePointRangeScore int           `json:"expirePointRangeScore"`
	CardNum               int           `json:"cardNum"`
}
