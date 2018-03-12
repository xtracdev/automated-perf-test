import { TestBed, inject } from '@angular/core/testing';

import { TestSuiteService } from './test-suite.service';

describe('TestSuiteService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [TestSuiteService]
    });
  });

  it('should be created', inject([TestSuiteService], (service: TestSuiteService) => {
    expect(service).toBeTruthy();
  }));
});
