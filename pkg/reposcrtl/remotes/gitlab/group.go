package gitlab

import (
	"fmt"
	"strconv"

	generic "github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
	log "github.com/sirupsen/logrus"
	gl "github.com/xanzy/go-gitlab"
)

type fetchGroupProjects struct {
	client *gl.Client
}

func (p fetchGroupProjects) Get(id string) ([]generic.RemoteRepoInfo, error) {
	log.Debugf("Load Projects from Gitlab Group %s", id)

	if groupID, err := strconv.Atoi(id); err == nil {
		group, _, err := p.client.Groups.GetGroup(groupID, nil, nil)
		if err != nil {
			return nil, err
		}

		log.Debugf("Group %s with %d projects", group.Name, len(group.Projects))

		var elements []generic.RemoteRepoInfo

		groups, _, err := p.client.Groups.ListSubgroups(groupID, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, element := range groups {
			subProjects, err := p.Get(strconv.Itoa(element.ID))
			if err != nil {
				return nil, err
			}

			elements = append(elements, subProjects...)
		}

		for _, element := range group.Projects {
			elements = append(elements, &Project{
				ref: element,
			})
		}

		return elements, nil
	}

	return nil, fmt.Errorf(fmt.Sprintf("Group with %s not found", id))
}
