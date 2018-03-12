import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TestCasesListComponent } from './test-cases-list.component';

describe('TestCasesListComponent', () => {
  let component: TestCasesListComponent;
  let fixture: ComponentFixture<TestCasesListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TestCasesListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TestCasesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
