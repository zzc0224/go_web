package redis

import "strconv"

func DeletePost(PostId uint64, communityName string) {
	client.ZRem(KeyPostScoreZSet, PostId)
	client.ZRem(KeyPostTimeZSet, PostId)
	key := KeyCommunityPostSetPrefix + communityName
	client.SRem(key, PostId)
	key = KeyPostInfoHashPrefix + strconv.FormatUint(PostId, 10)
	client.Del(key)
	key = KeyPostVotedZSetPrefix + strconv.FormatUint(PostId, 10)
	client.Del(key)
}
