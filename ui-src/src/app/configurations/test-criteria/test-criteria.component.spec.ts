import { async, ComponentFixture, TestBed } from "@angular/core/testing";

import { TestCriteriaComponent } from "./test-criteria.component";

describe("TestCriteriaComponent", () => {
  let component: TestCriteriaComponent;
  let fixture: ComponentFixture<TestCriteriaComponent>;

  beforeEach(
    async(() => {
      TestBed.configureTestingModule({
        declarations: [TestCriteriaComponent]
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TestCriteriaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
