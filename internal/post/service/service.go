package service

type Service struct {
	Post    *PostService
	Tag     *TagService
	Like    *LikeService
	Comment *CommentService
}

func New(postRepo PostRepo, tagRepo TagRepo, likeRepo LikeRepo, commentRepo CommentRepo) *Service {
	ts := NewTagService(tagRepo)
	ps := NewPostService(postRepo, ts)
	ls := NewLikeService(likeRepo)
	cs := NewCommentService(commentRepo)

	return &Service{
		Post:    ps,
		Tag:     ts,
		Like:    ls,
		Comment: cs,
	}
}
