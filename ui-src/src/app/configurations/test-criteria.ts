export class TestCriteria {
  constructor(
    public numIterations: number,
    public concurrentUsers: number,
    public memoryVariance: number,
    public serviceVariance: number,
    public requestDelay: number,
    public tpsFreq: number,
    public rampUsers: number,
    public rampDelay: number,
    public testSuite: string
  ) {}
}
