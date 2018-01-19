import { TestBed, inject } from "@angular/core/testing";
import { HttpClientModule } from "@angular/common/http";
import { AutomatedUIServices } from "./automated-ui-services";
import {
  HttpClientTestingModule,
  HttpTestingController
} from "@angular/common/http/testing";
import { MOCKDATA } from "./mockData";

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
  // it(
  //   "should be defined",
  //   inject(
  //     [AutomatedUIServices, HttpTestingController],
  //     (service: AutomatedUIServices, httpMock: HttpTestingController) => {
  //       service.createJsonFile(configData).subscribe(logs => {
  //         expect(spyOn(logs)).toEqual("Created JSON file");
  //       });
  //     }
  //   )
  // );
  it(
    "should be defined",
    inject(
      [AutomatedUIServices, HttpTestingController],
      (service: AutomatedUIServices, httpMock: HttpTestingController) => {
        service.createJsonFile(MOCKDATA).subscribe(res => {
          expect(res).toBeUndefined();
        });
      }
    )
  );
  // it(
  //   "should be defined",
  //   inject(
  //     [AutomatedUIServices, HttpTestingController],
  //     (service: AutomatedUIServices, httpMock: HttpTestingController) => {
  //       service.createJsonFile(configData).subscribe(res => {
  //         expect(logs).toBeUndefined();
  //       });
  //     }
  //   )
  // );
});
