package utils

var K3sUtil = newDockerUtil()

type k3sUtil struct {
	//client *helmclient.HelmClient
}

//func (u *k3sUtil) ListRelease(_ context.Context) (releases []*release.Release, err error) {
//	releases, err = u.client.ListDeployedReleases()
//	if err != nil {
//		return nil, err
//	}
//	return releases, nil
//}
//
//func (u k3sUtil) Install(ctx context.Context, valuesYaml string, releaseName string) (release *release.Release, err error) {
//	releaseInfo, err := u.client.InstallOrUpgradeChart(ctx, &helmclient.ChartSpec{
//		ReleaseName: releaseName,
//		ChartName:   g.Cfg().GetString("123"),
//		Namespace:   "default",
//		UpgradeCRDs: true,
//		Wait:        true,
//		ValuesYaml:  valuesYaml,
//	})
//	if err != nil {
//		return nil, err
//	}
//	return releaseInfo, nil
//}
