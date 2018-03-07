import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {TestCasesComponent} from "./test-cases.component";
import {Observable} from "rxjs/Observable";
import {environment} from "../../environments/environment.prod";

@Injectable()
export class TestCaseService {

  private headers: HttpHeaders;
  constructor(private http: HttpClient) {
    this.headers = new HttpHeaders({"Content-Type": "application/json"});
  }

  getAllCases$(testCasesPath): Observable<any> {
    this.headers = this.headers.set("testSuitePathDir", testCasesPath);
    return this.http.get(`${environment.apiBaseUrl}test-suites` + '/getAllCases/', {
      headers: this.headers
    });
    

    
  }
  postTestCases$(testCaseData, testCasePath): Observable<any> {
    this.headers = this.headers.set("testCasePathDir", testCasePath);
    return this.http.post(`${environment.apiBaseUrl}test-cases`, testCaseData, {
      headers: this.headers
    });
  }
  getOneTestCase$(testCasePath, testCaseFileName): Observable<any> {
    this.headers = this.headers.set("testCasePathDir", testCasePath);
    return this.http.get(`${environment.apiBaseUrl}test-cases/${testCaseFileName}`, {
      headers: this.headers
    });
  }

  putTestCase$(testCaseData, testCasePath, testCaseFileName): Observable<any> {
    this.headers = this.headers.set("testCasePathDir", testCasePath);
    return this.http.put(
      `${environment.apiBaseUrl}test-cases/${testCaseFileName}`,
      testCaseData,
      { headers: this.headers }
    );
  }

}
