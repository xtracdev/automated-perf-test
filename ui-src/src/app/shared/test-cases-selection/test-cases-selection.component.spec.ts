import { async, ComponentFixture, TestBed } from "@angular/core/testing";
import {TestCaseService} from "../../test-cases/test-case.service";
import { TestCasesSelectionComponent } from "./test-cases-selection.component";
import {CUSTOM_ELEMENTS_SCHEMA, NO_ERRORS_SCHEMA} from "@angular/core";
import {HttpClientModule} from "@angular/common/http";

describe("TestCasesSelectionComponent", () => {
  let component: TestCasesSelectionComponent;
  let fixture: ComponentFixture<TestCasesSelectionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TestCasesSelectionComponent ],
      providers: [TestCaseService],
      imports: [HttpClientModule],
      schemas: [NO_ERRORS_SCHEMA,
        CUSTOM_ELEMENTS_SCHEMA]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TestCasesSelectionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
