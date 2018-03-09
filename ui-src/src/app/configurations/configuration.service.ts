import {Injectable} from "@angular/core";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {ConfigurationsComponent} from "./configurations.component";
import {Observable} from "rxjs/Observable";
import {environment} from "../../environments/environment.prod";


@Injectable()
export class ConfigurationService {
  private headers: HttpHeaders;
  constructor(private http: HttpClient) {
    this.headers = new HttpHeaders({"Content-Type": "application/json"});
  }

  postConfig$(configData, configPath): Observable<any> {
    this.headers = this.headers.set("configPathDir", configPath);
    return this.http.post(`${environment.apiBaseUrl}configs`, configData, {
      headers: this.headers
    });


  }
  
  getConfig$(configPath, xmlFileName): Observable<any> {
    this.headers = this.headers.set("configPathDir", configPath);
    return this.http.get(`${environment.apiBaseUrl}configs/${xmlFileName}`, {
      headers: this.headers
    });
  }

  putConfig$(configData, configPath, xmlFileName): Observable<any> {
    this.headers = this.headers.set("configPathDir", configPath);
    return this.http.put(
      `${environment.apiBaseUrl}configs/${xmlFileName}`,
      configData,
      {headers: this.headers}
    );
  }

  getSchema$(location: string): Observable<any> {
    return this.http
      .get(`http://localhost:4200/${location}`, {headers: this.headers})
      .map((data: any) => data);
  }

  getAllCases$(testCasePath): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testCasePath);
    return this.http.get(`${environment.apiBaseUrl}test-cases`, {
      headers: this.headers
    });
  }

  postTestCases$(testCaseData, testCasePath): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testCasePath);
    return this.http.post(`${environment.apiBaseUrl}test-cases`, testCaseData, {
      headers: this.headers
    });
  }

  getOneTestCase$(testCasePath, testCaseFileName): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testCasePath);
    return this.http.get(`${environment.apiBaseUrl}test-cases/${testCaseFileName}`, {
      headers: this.headers
    });
  }

  putTestCase$(testCaseData, testCasePath, testCaseFileName): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testCasePath);
    return this.http.put(
      `${environment.apiBaseUrl}test-cases/${testCaseFileName}`,
      testCaseData,
      { headers: this.headers }
    );
  }

}
