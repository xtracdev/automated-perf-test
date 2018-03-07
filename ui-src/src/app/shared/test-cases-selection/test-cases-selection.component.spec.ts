import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TestCasesSelectionComponent } from './test-cases-selection.component';

describe('TestCasesSelectionComponent', () => {
  let component: TestCasesSelectionComponent;
  let fixture: ComponentFixture<TestCasesSelectionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TestCasesSelectionComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TestCasesSelectionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
