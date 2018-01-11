export class TestCriteria {
  public numIterations: number;
  public concurrentUsers: number;
  public memoryVariance: number;
  public serviceVariance: number;
  public requestDelay: number;
  public tpsFreq: number;
  public rampUsers: number;
  public rampDelay: number;
  public testSuite: string;

  constructor(
    numIterations?: number,
    concurrentUsers?: number,
    memoryVariance?: number,
    serviceVariance?: number,
    requestDelay?: number,
    tpsFreq?: number,
    rampUsers?: number,
    rampDelay?: number,
    testSuite?: string
  ) {
    this.numIterations = numIterations;
    this.concurrentUsers = concurrentUsers;
    this.memoryVariance = memoryVariance;
    this.serviceVariance = serviceVariance;
    this.requestDelay = requestDelay;
    this.tpsFreq = tpsFreq;
    this.rampUsers = rampUsers;
    this.rampDelay = rampDelay;
    this.testSuite = testSuite;
  }
}
