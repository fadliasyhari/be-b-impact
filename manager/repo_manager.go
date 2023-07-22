package manager

import (
	"be-b-impact.com/csr/repository"
	firebase "firebase.google.com/go/v4"
)

// RepositoryManager -> all repo
type RepositoryManager interface {
	UsersRepo() repository.UsersRepository
	CategoryRepo() repository.CategoryRepository
	TagRepo() repository.TagRepository
	ContentRepo() repository.ContentRepository
	ImageRepo() repository.ImageRepository
	TagsContentRepo() repository.TagsContentRepository
	ProposalRepo() repository.ProposalRepository
	ProposalDetailRepo() repository.ProposalDetailRepository
	FileRepo() repository.FileRepository
	ProgressRepo() repository.ProgressRepository
	ProposalProgressRepo() repository.ProposalProgressRepository
	EventRepo() repository.EventRepository
	EventParticipantRepo() repository.EventParticipantRepository
	EventImageRepo() repository.EventImageRepository
	NotificationRepo() repository.NotificationRepository
	UserDetailRepo() repository.UserDetailRepository
}

type repositoryManager struct {
	infra       InfraManager
	firebaseApp *firebase.App
}

// UserDetailRepo implements RepositoryManager.
func (r *repositoryManager) UserDetailRepo() repository.UserDetailRepository {
	return repository.NewUserDetailRepository(r.infra.Conn(), r.firebaseApp)
}

// NotificationRepo implements RepositoryManager.
func (r *repositoryManager) NotificationRepo() repository.NotificationRepository {
	return repository.NewNotificationRepository(r.infra.Conn())
}

// EventImageRepo implements RepositoryManager.
func (r *repositoryManager) EventImageRepo() repository.EventImageRepository {
	return repository.NewEventImageRepository(r.infra.Conn(), r.firebaseApp)
}

// EventParticipantRepo implements RepositoryManager.
func (r *repositoryManager) EventParticipantRepo() repository.EventParticipantRepository {
	return repository.NewEventParticipantRepository(r.infra.Conn())
}

// EventRepo implements RepositoryManager.
func (r *repositoryManager) EventRepo() repository.EventRepository {
	return repository.NewEventRepository(r.infra.Conn())
}

// ProgressRepo implements RepositoryManager.
func (r *repositoryManager) ProgressRepo() repository.ProgressRepository {
	return repository.NewProgressRepository(r.infra.Conn())
}

// ProposalProgressRepo implements RepositoryManager.
func (r *repositoryManager) ProposalProgressRepo() repository.ProposalProgressRepository {
	return repository.NewProposalProgressRepository(r.infra.Conn())
}

// FileRepo implements RepositoryManager.
func (r *repositoryManager) FileRepo() repository.FileRepository {
	return repository.NewFileRepository(r.infra.Conn(), r.firebaseApp)
}

// ProposalDetailRepo implements RepositoryManager.
func (r *repositoryManager) ProposalDetailRepo() repository.ProposalDetailRepository {
	return repository.NewProposalDetailRepository(r.infra.Conn())
}

// ProposalRepo implements RepositoryManager.
func (r *repositoryManager) ProposalRepo() repository.ProposalRepository {
	return repository.NewProposalRepository(r.infra.Conn())
}

// TagsContentRepo implements RepositoryManager.
func (r *repositoryManager) TagsContentRepo() repository.TagsContentRepository {
	return repository.NewTagsContentRepository(r.infra.Conn())
}

// ImageRepo implements RepositoryManager
func (r *repositoryManager) ImageRepo() repository.ImageRepository {
	return repository.NewImageRepository(r.infra.Conn(), r.firebaseApp)
}

// ContentRepo implements RepositoryManager
func (r *repositoryManager) ContentRepo() repository.ContentRepository {
	return repository.NewContentRepository(r.infra.Conn())
}

// TagRepo implements RepositoryManager
func (r *repositoryManager) TagRepo() repository.TagRepository {
	return repository.NewTagRepository(r.infra.Conn())
}

// CategoryRepo implements RepositoryManager
func (r *repositoryManager) CategoryRepo() repository.CategoryRepository {
	return repository.NewCategoryRepository(r.infra.Conn())
}

func (r *repositoryManager) UsersRepo() repository.UsersRepository {
	return repository.NewUsersRepository(r.infra.Conn())
}

func NewRepositoryManager(infra InfraManager) RepositoryManager {
	firebaseApp := infra.FirebaseApp()
	return &repositoryManager{
		infra:       infra,
		firebaseApp: firebaseApp,
	}
}
