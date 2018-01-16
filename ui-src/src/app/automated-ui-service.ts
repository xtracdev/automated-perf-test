import { Injectable } from "@angular/core";
import { HttpModule } from "@angular/http";
import { Http } from "@angular/http";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { ConfigurationsComponent } from "./configurations/configurations.component";
import { Headers } from "@angular/http/src/headers";

const httpOptions = {
  headers: new HttpHeaders({
    "Content-Type": "application/json",
    configPathDir:
      "C:/Users/a615194/go/src/github.com/xtracdev/automated-perf-test/config/"
    //this will be changed to whatever the person has selected for
    // their config file path in the form.
  })
};

@Injectable()
export class AutomatedUIService {
  constructor(private http: HttpClient) {}

  private url = "http://localhost:9191/configs";

  addConfig(form: Data): void {
    console.log("Form", form);
    this.http.post(this.url, form, httpOptions).subscribe(data => {
      console.log(data);
    });
  }
}

export class Data {
  apiName: string;
  targetHost: string;
  targetPort: number;
  numIterations: number;
  concurrentUsers: number;
  allowablePeakMemoryVariance: number;
  allowableServiceResponseTimeVariance: number;
  testSuite: string;
  requestDelay: number;
  TPSFreq: number;
  rampUsers: number;
  rampDelay: number;
  testCaseDir: string;
  testSuiteDir: string;
  baseStatsOutputDir: string;
  reportOutputDir: string;
}
