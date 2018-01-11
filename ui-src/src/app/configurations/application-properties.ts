export class ApplicationProperties {
  public applicationName: string;
  public targetHost: string;
  public targetPort: number;
  public memoryEndPoint: string;

  constructor(
    applicationName?: string,
    targetHost?: string,
    targetPort?: number,
    memoryEndPoint?: string
  ) {
    this.applicationName = applicationName;
    this.targetHost = targetHost;
    this.targetPort = targetPort;
    this.memoryEndPoint = memoryEndPoint;
  }
}
