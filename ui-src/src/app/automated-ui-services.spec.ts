import { TestBed, inject } from "@angular/core/testing";
import { HttpClientModule } from "@angular/common/http";
import { AutomatedUIServices } from "./automated-ui-services";
import {
  HttpClientTestingModule,
  HttpTestingController
} from "@angular/common/http/testing";

describe("AutomatedUIServices", () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientModule, HttpClientTestingModule],
      providers: [AutomatedUIServices]
    });
  });

  it(
    "should be created",
    inject([AutomatedUIServices], (service: AutomatedUIServices) => {
      expect(service).toBeTruthy();
    })
  );
});
