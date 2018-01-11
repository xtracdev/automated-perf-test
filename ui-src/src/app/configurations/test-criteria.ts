export class TestCriteria {
  constructor(
    public numIterationsTxt: number,
    public concurrentUsersTxt: number,
    public memoryVarianceTxt: number,
    public serviceVarianceTxt: number,
    public requestDelayTxt: number,
    public tpsFreqTxt: number,
    public rampUsersTxt: number,
    public rampDelayTxt: number,
    public testSuiteTxt: string
  ) {}
}
