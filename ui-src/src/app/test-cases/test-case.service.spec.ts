import { TestBed, inject } from "@angular/core/testing";

import { TestCaseService } from "./test-case.service";

describe("TestCaseService", () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [TestCaseService]
    });
  });

  it("should be created", inject([TestCaseService], (service: TestCaseService) => {
    expect(service).toBeTruthy();
  }));
});
