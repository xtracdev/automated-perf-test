import {async, fakeAsync, ComponentFixture, TestBed} from "@angular/core/testing";
import {TestCasesComponent} from "./test-cases.component";
import {TestSuiteService} from "../test-suites/test-suite.service";
import {ConfigurationService} from "../configurations/configuration.service";
import {TestCaseService} from "../test-cases/test-case.service";
import {CUSTOM_ELEMENTS_SCHEMA, NO_ERRORS_SCHEMA} from "@angular/core";
import {ToastModule} from "ng2-toastr/ng2-toastr";
import {HttpClientModule} from "@angular/common/http";
import {FormsModule} from "@angular/forms";
import {ToastsManager, ToastOptions} from "ng2-toastr/ng2-toastr";
import {Observable} from "rxjs/Observable";
import "rxjs/add/observable/of";
import "rxjs/add/observable/throw";

describe("TestCasesComponent", () => {
  let component: TestCasesComponent;
  let fixture: ComponentFixture<TestCasesComponent>;
  let testSuiteService, toastr;

  beforeEach(
    async(() => {
      TestBed.configureTestingModule({
        providers: [ConfigurationService, ToastsManager, ToastOptions, TestSuiteService, TestCaseService],
        declarations: [TestCasesComponent],
        imports: [FormsModule, HttpClientModule, ToastModule.forRoot()],
        schemas: [NO_ERRORS_SCHEMA,
          CUSTOM_ELEMENTS_SCHEMA]
      }).compileComponents();

      fixture = TestBed.createComponent(TestCasesComponent);
      component = fixture.componentInstance;
      fixture.detectChanges();

      testSuiteService = TestBed.get(TestSuiteService);
      toastr = TestBed.get(ToastsManager);
    }));

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
