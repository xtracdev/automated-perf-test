export class OutputPaths {
  public testCaseDicectory: string;
  public testSuiteDirectory: string;
  public baseSuiteDirectory: string;
  public reportsDirectory: string;

  constructor(
    testCaseDicectory?: string,
    testSuiteDirectory?: string,
    baseSuiteDirectory?: string,
    reportsDirectory?: string
  ) {
    this.testCaseDicectory = testCaseDicectory;
    this.testSuiteDirectory = testSuiteDirectory;
    this.baseSuiteDirectory = baseSuiteDirectory;
    this.reportsDirectory = reportsDirectory;
  }
}
