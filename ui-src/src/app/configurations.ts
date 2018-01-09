export class Configurations {
  constructor(
    public configPathDir: string,
    public applicationNameTxt: string,
    public targetHostTxt: string,
    public targetPortTxt: number,
    public memoryEndPointTxt: string,
    public numIterationsTxt: number,
    public concurrentUsersTxt: number,
    public memoryVarianceTxt: number,
    public serviceVarianceTxt: number,
    public testSuiteTxt: string,
    public requestDelayTxt: number,
    public tpsFreqTxt: number,
    public rampUsersTxt: number,
    public rampDelayTxt: number,
    public testCaseDirTxt: string,
    public testSuiteDirTxt: string,
    public baseStatsDirTxt: string,
    public reportsDirTxt: string
  ) {}
}
