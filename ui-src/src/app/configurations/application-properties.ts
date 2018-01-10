export class ApplicationProperties {
  constructor(
    public applicationName: string,
    public targetHost: string,
    public targetPort: number,
    public memoryEndPoint: string
  ) {}
}
