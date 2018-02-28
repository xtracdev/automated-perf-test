import {TestBed, inject, getTestBed} from '@angular/core/testing';
import {HttpClientModule} from "@angular/common/http";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {environment} from "../../environments/environment.prod";
import {
  HttpClientTestingModule,
  HttpTestingController
} from "@angular/common/http/testing";

import {ConfigurationService} from './configuration.service';

describe('ConfigurationService', () => {
  let injector;
  let service;
  let httpInterceptor: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientModule, HttpClientTestingModule],
      providers: [ConfigurationService]
    });

    injector = getTestBed();
    service = injector.get(ConfigurationService);
    httpInterceptor = injector.get(HttpTestingController);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });

  it("should make post config request", () => {
    // TODO add mock so it is no longer calling the real service
    service.postConfig$({data: "data1"}, "data2").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "data2");

    const req = httpInterceptor.expectOne(`${environment.apiBaseUrl}configs`);
    expect(req.request.method).toBe("POST");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toEqual({data: "data1"});
    req.flush({});
  });

  it("should make put config request", () => {
    // TODO add mock so it is no longer calling the real service
    service.putConfig$({data: "data1"}, "direct path", "fileName.xml").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "direct path");

    const req = httpInterceptor.expectOne(`${environment.apiBaseUrl}configs/fileName.xml`);
    expect(req.request.method).toBe("PUT");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toEqual({data: "data1"});
    req.flush({});
  });

  it("should make get config request", () => {
    // TODO add mock so it is no longer calling the real service
    service.getConfig$("direct path", "fileName.xml").subscribe(() => {});

    let headers = new HttpHeaders();
    headers = headers.set("Content-Type", "application/json;");
    headers = headers.append("configPathDir", "direct path");

    const req = httpInterceptor.expectOne(`${environment.apiBaseUrl}configs/fileName.xml`);
    expect(req.request.method).toBe("GET");
    expect(req.request.headers.getAll).toBe(headers.getAll);
    expect(req.request.body).toBeNull();
    req.flush({});
  });
});
