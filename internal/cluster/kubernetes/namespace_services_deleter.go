// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"context"

	"emperror.dev/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/banzaicloud/pipeline/pkg/k8sclient"
)

type NamespaceServicesDeleter struct {
	logger logrus.FieldLogger
}

func MakeNamespaceServicesDeleter(logger logrus.FieldLogger) NamespaceServicesDeleter {
	return NamespaceServicesDeleter{
		logger: logger,
	}
}

func (d NamespaceServicesDeleter) Delete(ctx context.Context, organizationID uint, clusterName string, k8sConfig []byte, namespace string) error {
	logger := d.logger.WithField("organizationID", organizationID).WithField("clusterName", clusterName).WithField("namespace", namespace)

	client, err := k8sclient.NewClientFromKubeConfig(k8sConfig)
	if err != nil {
		return err
	}
	services, err := client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return errors.WrapIf(err, "could not list services to delete")
	}

	for _, service := range services.Items {
		switch service.Name {
		case "kubernetes":
			continue
		}
		err := retry(func() error {
			logger.Infof("deleting kubernetes service %q", service.Name)
			err := client.CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{})
			if err != nil {
				return errors.WrapIff(err, "failed to delete %q service", service.Name)
			}
			return nil
		}, 3, 1)
		if err != nil {
			return err
		}
	}
	err = retry(func() error {
		services, err := client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return errors.WrapIf(err, "could not list remaining services")
		}
		left := []string{}
		for _, svc := range services.Items {
			switch svc.Name {
			case "kubernetes":
				continue
			default:
				logger.Infof("service %q still %v", svc.Name, svc.Status)
				left = append(left, svc.Name)
			}
		}
		if len(left) > 0 {
			return errors.WithDetails(errors.Errorf("services remained after deletion: %v", left), "services", left)
		}
		return nil
	}, 6, 30)
	return err
}
