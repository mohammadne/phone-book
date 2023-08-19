deploy:
	helmsman -apply -f ./deployments/helmsman.yaml

undeploy:
	helmsman -destroy -f ./deployments/helmsman.yaml
