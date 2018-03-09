import {TestBed, inject, getTestBed} from "@angular/core/testing";
import {HttpClientModule} from "@angular/common/http";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {
  HttpClientTestingModule,
  HttpTestingController
} from "@angular/common/http/testing";

import { TestCaseService } from "./test-case.service";

describe("TestCaseService", () => {
  let injector;
  let service;
  let httpInterceptor: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientModule, HttpClientTestingModule],
      providers: [TestCaseService]
    });

    injector = getTestBed();
    service = injector.get(TestCaseService);
    httpInterceptor = injector.get(HttpTestingController);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });

  it("should make post test case request", () => {
    // TODO add mock so it is no longer calling the real service
    service.postTestCase$({data: "data1"}, "data2").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "data2");

    const req = httpInterceptor.expectOne(`${environment.API_BASE_URL}test-cases`);
    expect(req.request.method).toBe("POST");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toEqual({data: "data1"});
    req.flush({});
  });

  it("should make put test case request", () => {
    // TODO add mock so it is no longer calling the real service
    service.putTestCase$({data: "data1"}, "direct path", "fileName.xml").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "direct path");

    const req = httpInterceptor.expectOne(`${environment.API_BASE_URL}test-cases/fileName.xml`);
    expect(req.request.method).toBe("PUT");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toEqual({data: "data1"});
    req.flush({});
  });

  it("should make get test case request", () => {
    // TODO add mock so it is no longer calling the real service
    service.getTestCase$("direct path", "fileName.xml").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "direct path");

    const req = httpInterceptor.expectOne(`${environment.API_BASE_URL}test-cases/fileName.xml`);
    expect(req.request.method).toBe("GET");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toBeNull();
    req.flush({});
  });

  it("should make get all test case request", () => {
    // TODO add mock so it is no longer calling the real service
    service.getTestCases$("direct path", "fileName.xml").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "direct path");

    const req = httpInterceptor.expectOne(`${environment.API_BASE_URL}test-cases`);
    expect(req.request.method).toBe("GET");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toBeNull();
    req.flush({});
  });
});
