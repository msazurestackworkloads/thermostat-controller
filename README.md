## a basic kubernetes thermostat controller/operator

creates and reconciles CRs on a cluster provided a kube api server url

attempts to match current temperature to desired temperature by adjusting the current temperature by 1 each reconcile loop and deletes the CR once complete
