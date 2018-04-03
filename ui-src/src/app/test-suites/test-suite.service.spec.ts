import {TestBed, inject, getTestBed} from "@angular/core/testing";
import {HttpClientModule} from "@angular/common/http";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {
  HttpClientTestingModule,
  HttpTestingController
} from "@angular/common/http/testing";

import {TestSuiteService} from "./test-suite.service";

describe("TestSuiteService", () => {
  let injector;
  let service;
  let httpInterceptor: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientModule, HttpClientTestingModule],
      providers: [TestSuiteService]
    });

    injector = getTestBed();
    service = injector.get(TestSuiteService);
    httpInterceptor = injector.get(HttpTestingController);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });

  it("should make post test suite request", () => {
    // TODO add mock so it is no longer calling the real service
    service.postTestSuite$({data: "data1"}, "data2").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "data2");

    const req = httpInterceptor.expectOne(`${environment.API_BASE_URL}test-suites`);
    expect(req.request.method).toBe("POST");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toEqual({data: "data1"});
    req.flush({});
  });

  it("should make put test suite request", () => {
    // TODO add mock so it is no longer calling the real service
    service.putTestSuite$({data: "data1"}, "direct path", "fileName.xml").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "direct path");

    const req = httpInterceptor.expectOne(`${environment.API_BASE_URL}test-suites/fileName.xml`);
    expect(req.request.method).toBe("PUT");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toEqual({data: "data1"});
    req.flush({});
  });

  it("should make get test suite request", () => {
    // TODO add mock so it is no longer calling the real service
    service.getTestSuite$("direct path", "fileName.xml").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "direct path");

    const req = httpInterceptor.expectOne(`${environment.API_BASE_URL}test-suites/fileName.xml`);
    expect(req.request.method).toBe("GET");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toBeNull();
    req.flush({});
  });

  it("should make get all test suites request", () => {
    // TODO add mock so it is no longer calling the real service
    service.getTestSuites$("direct path", "fileName.xml").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "direct path");

    const req = httpInterceptor.expectOne(`${environment.API_BASE_URL}test-suites`);
    expect(req.request.method).toBe("GET");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toBeNull();
    req.flush({});
  });
});
