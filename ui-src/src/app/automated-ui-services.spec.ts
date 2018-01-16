import { TestBed, inject } from "@angular/core/testing";

import { AutomatedUIServices } from "./automated-ui-services";

describe("PostService", () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
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
