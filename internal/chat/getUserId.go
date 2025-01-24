package chat

func (chat *Chat) GetUserId() uuId {
	return chat.Client.UserId
}
