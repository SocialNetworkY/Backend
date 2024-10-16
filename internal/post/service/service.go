package service

type Service struct {
	Post    *PostService
	Tag     *TagService
	Like    *LikeService
	Comment *CommentService
}

func New(postRepo PostRepo, tagRepo TagRepo, likeRepo LikeRepo, commentRepo CommentRepo, rg ReportGateway) *Service {
	ts := NewTagService(tagRepo)
	ls := NewLikeService(likeRepo)
	cs := NewCommentService(commentRepo)
	ps := NewPostService(postRepo, rg, ts, cs, ls)

	return &Service{
		Post:    ps,
		Tag:     ts,
		Like:    ls,
		Comment: cs,
	}
}
