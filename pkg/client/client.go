package client

import (
	v1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client represents the Chart Manager client.
type Client struct {
	Clientset              *clientset.Clientset
	RESTClient             *rest.RESTClient
	APIExtensionsClientset *apiextensionsclientset.Clientset
}

// NewForConfig instantiates and returns the client and scheme.
func NewForConfig(cfg *rest.Config) (*Client, *runtime.Scheme, error) {
	s := runtime.NewScheme()
	c, err := initClients(cfg, s)
	if err != nil {
		return nil, nil, err
	}
	return c, s, nil
}

func initClients(cfg *rest.Config, s *runtime.Scheme) (*Client, error) {
	client, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	restconfig := restConfig(cfg, s)
	restclient, err := rest.RESTClientFor(&restconfig)
	if err != nil {
		return nil, err
	}

	// Instantiate the Kubernetes API extensions client.
	apiextensionsclient, err := apiextensionsclientset.NewForConfig(&restconfig)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Clientset:              client,
		RESTClient:             restclient,
		APIExtensionsClientset: apiextensionsclient,
	}
	return c, nil
}

func restConfig(cfg *rest.Config, s *runtime.Scheme) rest.Config {
	config := *cfg
	config.GroupVersion = &v1alpha1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(s)}
	return config
}
