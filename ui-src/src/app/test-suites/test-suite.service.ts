import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {TestSuitesComponent} from "./test-suites.component";
import {Observable} from "rxjs/Observable";
import {environment} from "../../environments/environment.prod";

@Injectable()
export class TestSuiteService {

  private headers: HttpHeaders;
  constructor(private http: HttpClient) {
    this.headers = new HttpHeaders({"Content-Type": "application/json"});
  }

  postTestSuite$(testSuiteData, testSuitePath): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testSuitePath);
    return this.http.post(`${environment.apiBaseUrl}test-suites`, testSuiteData, {
      headers: this.headers
    });
  }

  getTestSuite$(testSuitePath, testSuiteFileName): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testSuitePath);
    return this.http.get(`${environment.apiBaseUrl}test-suites/${testSuiteFileName}`, {
      headers: this.headers
    });
  }

  getAllTestSuite$(testSuitePath): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testSuitePath);
    return this.http.get(`${environment.apiBaseUrl}test-suites`, {
      headers: this.headers
    });
  }
  getSchema$(location: string): Observable<any> {
    return this.http
      .get(`http://localhost:4200/${location}`, {headers: this.headers})
      .map((data: any) => data);
  }

  putTestSuite$(testSuiteData, testSuitePath, testSuiteFileName): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testSuitePath);
    return this.http.put(
      `${environment.apiBaseUrl}test-suites/${testSuiteFileName}`,
      testSuiteData,
      {headers: this.headers}
    );
  }
}
