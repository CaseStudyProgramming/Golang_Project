# Component Documentation

This document provides comprehensive documentation for all components in the Task Manager application, organized by feature modules and shared utilities.

## Table of Contents

- [Shared Components](#shared-components)
- [Task Management Components](#task-management-components)
- [Category Components](#category-components)
- [Tag Components](#tag-components)
- [Analytics Components](#analytics-components)
- [Component Usage Best Practices](#component-usage-best-practices)

## Shared Components

Shared components are reusable UI primitives located in `src/lib/shared/components/`. These components can be used across any feature module.

### PasswordStrength

A password strength indicator component that provides visual feedback on password complexity.

**Location:** `src/lib/shared/components/PasswordStrength.svelte`

**Usage:**
```svelte
<script>
  import { PasswordStrength } from '$lib/shared/components';
  let password = $state('');
</script>

<PasswordStrength password={password} />
```

**Props:**
- `password` (string): The password string to evaluate

**Features:**
- Visual strength indicator (weak, medium, strong)
- Color-coded feedback
- Real-time strength evaluation

### Skeleton

A loading skeleton component for displaying placeholder content while data is loading.

**Location:** `src/lib/shared/components/Skeleton.svelte`

**Usage:**
```svelte
<script>
  import { Skeleton } from '$lib/shared/components';
</script>

<div class="space-y-4">
  <Skeleton class="h-4 w-3/4" />
  <Skeleton class="h-4 w-1/2" />
  <Skeleton class="h-32 w-full" />
</div>
```

**Props:**
- `class` (string): Additional CSS classes for styling

**Features:**
- Configurable size and shape
- Smooth loading animation
- Accessible placeholder content

### LoadingSpinner

A circular loading spinner for indicating loading states.

**Location:** `src/lib/shared/components/LoadingSpinner.svelte`

**Usage:**
```svelte
<script>
  import { LoadingSpinner } from '$lib/shared/components';
</script>

<LoadingSpinner size="medium" />
```

**Props:**
- `size` (string): Size variant ('small', 'medium', 'large')

**Features:**
- Multiple size options
- Smooth animation
- Accessible loading indicator

### ProgressBar

A progress bar component for displaying task completion or loading progress.

**Location:** `src/lib/shared/components/ProgressBar.svelte`

**Usage:**
```svelte
<script>
  import { ProgressBar } from '$lib/shared/components';
</script>

<ProgressBar progress={75} total={100} label="Task Progress" />
```

**Props:**
- `progress` (number): Current progress value
- `total` (number): Total value for progress calculation
- `label` (string): Optional label for the progress bar

**Features:**
- Visual progress indication
- Percentage display
- Configurable colors

### EmptyState

A component for displaying empty states when no data is available.

**Location:** `src/lib/shared/components/EmptyState.svelte`

**Usage:**
```svelte
<script>
  import { EmptyState } from '$lib/shared/components';
</script>

<EmptyState 
  title="No tasks found"
  description="Create your first task to get started"
  icon="📝"
/>
```

**Props:**
- `title` (string): Title for the empty state
- `description` (string): Descriptive text
- `icon` (string): Icon or emoji to display

**Features:**
- Customizable messaging
- Icon support
- Action button support

### Toast & ToastContainer

Toast notification system for displaying temporary messages.

**Location:** `src/lib/shared/components/Toast.svelte`, `src/lib/shared/components/ToastContainer.svelte`

**Usage:**
```svelte
<script>
  import { toastStore } from '$lib/shared/stores';
</script>

<!-- In your component logic -->
toastStore.success('Task created successfully');
toastStore.error('Failed to create task');
toastStore.info('Processing request...');

<!-- In your layout -->
<ToastContainer />
```

**Toast Store Methods:**
- `success(message: string)`: Display success toast
- `error(message: string)`: Display error toast
- `info(message: string)`: Display info toast
- `warning(message: string)`: Display warning toast

**Features:**
- Auto-dismiss after timeout
- Multiple toast types
- Stacked notifications
- Smooth animations

### ConfirmDialog & ConfirmDialogContainer

Confirmation dialog system for destructive actions.

**Location:** `src/lib/shared/components/ConfirmDialog.svelte`, `src/lib/shared/components/ConfirmDialogContainer.svelte`

**Usage:**
```svelte
<script>
  import { confirmStore } from '$lib/shared/stores';
  
  async function handleDelete() {
    const confirmed = await confirmStore.showConfirm({
      title: 'Delete Task',
      message: 'Are you sure you want to delete this task?',
      confirmText: 'Delete',
      cancelText: 'Cancel'
    });
    
    if (confirmed) {
      // Perform delete action
    }
  }
</script>

<!-- In your layout -->
<ConfirmDialogContainer />
```

**Confirm Store Methods:**
- `showConfirm(options: ConfirmOptions)`: Show confirmation dialog
  - `title` (string): Dialog title
  - `message` (string): Confirmation message
  - `confirmText` (string): Confirm button text
  - `cancelText` (string): Cancel button text

**Features:**
- Promise-based API
- Customizable messaging
- Keyboard support (Esc to cancel)
- Accessible dialog

## Task Management Components

Task management components are located in `src/lib/features/tasks/components/`.

### TaskForm

A comprehensive form for creating and editing tasks.

**Location:** `src/lib/features/tasks/components/TaskForm.svelte`

**Usage:**
```svelte
<script>
  import { TaskForm } from '$lib/features/tasks';
  
  function handleSubmit(data) {
    console.log('Task data:', data);
  }
</script>

<TaskForm 
  mode="create"
  onSubmit={handleSubmit}
  onCancel={() => console.log('Cancelled')}
/>
```

**Props:**
- `mode` ('create' | 'edit'): Form mode
- `initialData` (CreateTaskPayload | UpdateTaskPayload): Initial form data
- `onSubmit` (function): Callback when form is submitted
- `onCancel` (function): Callback when form is cancelled

**Features:**
- Title and description fields
- Priority selection (LOW, MEDIUM, HIGH, URGENT)
- Due date picker
- Category selection
- Tag management with TagInput
- Form validation with Zod
- Error handling and display

### TaskList

A list component for displaying tasks with actions.

**Location:** `src/lib/features/tasks/components/TaskList.svelte`

**Usage:**
```svelte
<script>
  import { TaskList } from '$lib/features/tasks';
  import { taskStore } from '$lib/features/tasks/stores/task.store';
  
  function handleViewTask(task) {
    console.log('View task:', task);
  }
  
  function handleEditTask(task) {
    console.log('Edit task:', task);
  }
  
  function handleDeleteTask(task) {
    console.log('Delete task:', task);
  }
</script>

<TaskList 
  tasks={taskStore.state.tasks}
  isLoading={taskStore.state.isLoading}
  onViewTask={handleViewTask}
  onEditTask={handleEditTask}
  onDeleteTask={handleDeleteTask}
/>
```

**Props:**
- `tasks` (Task[]): Array of tasks to display
- `isLoading` (boolean): Loading state
- `onViewTask` (function): Callback when viewing a task
- `onEditTask` (function): Callback when editing a task
- `onDeleteTask` (function): Callback when deleting a task

**Features:**
- Task cards with priority indicators
- Status badges
- Category and tag display
- Progress bars for task completion
- Action buttons (view, edit, delete)
- Responsive grid layout
- Empty state handling

### TaskFilters

Filter component for task list filtering.

**Location:** `src/lib/features/tasks/components/TaskFilters.svelte`

**Usage:**
```svelte
<script>
  import { TaskFilters } from '$lib/features/tasks';
  import { taskStore } from '$lib/features/tasks/stores/task.store';
</script>

<TaskFilters 
  filters={taskStore.state.filters}
  onFilterChange={(filters) => taskStore.setFilters(filters)}
/>
```

**Props:**
- `filters` (TaskFilters): Current filter state
- `onFilterChange` (function): Callback when filters change

**Features:**
- Status filter (all, active, completed)
- Priority filter
- Category filter
- Date range filter
- Reset filters button

### TaskSearch

Search component for finding tasks.

**Location:** `src/lib/features/tasks/components/TaskSearch.svelte`

**Usage:**
```svelte
<script>
  import { TaskSearch } from '$lib/features/tasks';
  import { taskStore } from '$lib/features/tasks/stores/task.store';
</script>

<TaskSearch 
  onSearch={(query) => taskStore.searchTasks(query)}
/>
```

**Props:**
- `onSearch` (function): Callback when search query changes

**Features:**
- Real-time search
- Debounced input
- Clear search button
- Search history

### Pagination

Pagination component for task lists.

**Location:** `src/lib/features/tasks/components/Pagination.svelte`

**Usage:**
```svelte
<script>
  import { Pagination } from '$lib/features/tasks';
  import { taskStore } from '$lib/features/tasks/stores/task.store';
</script>

<Pagination 
  currentPage={taskStore.state.pagination.page}
  totalPages={taskStore.state.pagination.totalPages}
  onPageChange={(page) => taskStore.setPage(page)}
/>
```

**Props:**
- `currentPage` (number): Current page number
- `totalPages` (number): Total number of pages
- `onPageChange` (function): Callback when page changes

**Features:**
- Page number buttons
- Previous/Next navigation
- Page size selector
- Jump to page functionality

### SubtaskList

Component for managing subtasks within a task.

**Location:** `src/lib/features/tasks/components/SubtaskList.svelte`

**Usage:**
```svelte
<script>
  import { SubtaskList } from '$lib/features/tasks';
</script>

<SubtaskList 
  subtasks={task.subtasks}
  onToggleSubtask={(id) => console.log('Toggle:', id)}
  onDeleteSubtask={(id) => console.log('Delete:', id)}
  onAddSubtask={(title) => console.log('Add:', title)}
/>
```

**Props:**
- `subtasks` (Subtask[]): Array of subtasks
- `onToggleSubtask` (function): Callback when toggling subtask
- `onDeleteSubtask` (function): Callback when deleting subtask
- `onAddSubtask` (function): Callback when adding subtask

**Features:**
- Checkbox for completion
- Subtask title display
- Add new subtask
- Delete subtask
- Progress calculation

### ActivityLog

Component for displaying task activity history.

**Location:** `src/lib/features/tasks/components/ActivityLog.svelte`

**Usage:**
```svelte
<script>
  import { ActivityLog } from '$lib/features/tasks';
</script>

<ActivityLog 
  activities={task.activities}
/>
```

**Props:**
- `activities` (Activity[]): Array of activity logs

**Features:**
- Timeline display
- Activity type icons
- Timestamp formatting
- User attribution
- Expandable details

### TaskListSkeleton

Loading skeleton for task list.

**Location:** `src/lib/features/tasks/components/TaskListSkeleton.svelte`

**Usage:**
```svelte
<script>
  import { TaskListSkeleton } from '$lib/features/tasks';
</script>

<TaskListSkeleton count={5} />
```

**Props:**
- `count` (number): Number of skeleton items to display

**Features:**
- Matches TaskList layout
- Smooth loading animation
- Configurable item count

## Category Components

Category management components located in `src/lib/features/categories/components/`.

### CategoryForm

Form for creating and editing categories.

**Location:** `src/lib/features/categories/components/CategoryForm.svelte`

**Usage:**
```svelte
<script>
  import { CategoryForm } from '$lib/features/categories';
  
  function handleSubmit(data) {
    console.log('Category data:', data);
  }
</script>

<CategoryForm 
  mode="create"
  onSubmit={handleSubmit}
  onCancel={() => console.log('Cancelled')}
/>
```

**Props:**
- `mode` ('create' | 'edit'): Form mode
- `initialData` (Category): Initial category data
- `onSubmit` (function): Callback when form is submitted
- `onCancel` (function): Callback when form is cancelled

**Features:**
- Category name input
- Color picker for category color
- Form validation
- Color preview

### CategoryList

List component for displaying categories.

**Location:** `src/lib/features/categories/components/CategoryList.svelte`

**Usage:**
```svelte
<script>
  import { CategoryList } from '$lib/features/categories';
  import { categoryStore } from '$lib/features/categories/stores/category.store';
</script>

<CategoryList 
  categories={categoryStore.state.categories}
  onEditCategory={(category) => console.log('Edit:', category)}
  onDeleteCategory={(category) => console.log('Delete:', category)}
/>
```

**Props:**
- `categories` (Category[]): Array of categories
- `onEditCategory` (function): Callback when editing category
- `onDeleteCategory` (function): Callback when deleting category

**Features:**
- Color-coded category cards
- Category name display
- Task count per category
- Edit and delete actions
- Grid layout

### CategoryListSkeleton

Loading skeleton for category list.

**Location:** `src/lib/features/categories/components/CategoryListSkeleton.svelte`

**Usage:**
```svelte
<script>
  import { CategoryListSkeleton } from '$lib/features/categories';
</script>

<CategoryListSkeleton count={4} />
```

**Props:**
- `count` (number): Number of skeleton items

**Features:**
- Matches CategoryList layout
- Smooth loading animation

## Tag Components

Tag management components located in `src/lib/features/tags/components/`.

### TagInput

Advanced tag input component with autocomplete and creation.

**Location:** `src/lib/features/tags/components/TagInput.svelte`

**Usage:**
```svelte
<script>
  import { TagInput } from '$lib/features/tags';
  import { tagStore } from '$lib/features/tags/stores/tag.store';
  
  let selectedTags = $state(['tag1', 'tag2']);
</script>

<TagInput 
  selectedTags={selectedTags}
  availableTags={tagStore.state.tags}
  onCreateTag={(name) => tagStore.createTag({ name })}
/>
```

**Props:**
- `selectedTags` (string[]): Array of selected tag IDs
- `availableTags` (Tag[]): Array of available tags
- `onCreateTag` (function): Callback when creating new tag

**Features:**
- Tag autocomplete
- Create new tags on-the-fly
- Remove tags with click
- Tag color display
- Keyboard navigation
- Duplicate prevention

## Analytics Components

Analytics and reporting components located in `src/lib/features/analytics/components/`.

### StatisticsCards

Display cards for key statistics.

**Location:** `src/lib/features/analytics/components/StatisticsCards.svelte`

**Usage:**
```svelte
<script>
  import { StatisticsCards } from '$lib/features/analytics';
</script>

<StatisticsCards 
  stats={{
    totalTasks: 100,
    completedTasks: 75,
    overdueTasks: 5,
    completionRate: 75
  }}
/>
```

**Props:**
- `stats` (Statistics): Statistics object with key metrics

**Features:**
- Key metric display
- Trend indicators
- Color-coded values
- Responsive grid

### CompletionRateChart

Chart showing task completion rate over time.

**Location:** `src/lib/features/analytics/components/CompletionRateChart.svelte`

**Usage:**
```svelte
<script>
  import { CompletionRateChart } from '$lib/features/analytics';
</script>

<CompletionRateChart 
  data={completionData}
  timePeriod="week"
/>
```

**Props:**
- `data` (ChartPoint[]): Array of data points
- `timePeriod` (string): Time period for chart

**Features:**
- Line chart visualization
- Time period selection
- Responsive sizing
- Interactive tooltips

### PriorityChart

Chart showing task distribution by priority.

**Location:** `src/lib/features/analytics/components/PriorityChart.svelte`

**Usage:**
```svelte
<script>
  import { PriorityChart } from '$lib/features/analytics';
</script>

<PriorityChart 
  data={priorityData}
/>
```

**Props:**
- `data` (PriorityData): Priority distribution data

**Features:**
- Pie or bar chart
- Color-coded by priority
- Percentage display
- Interactive legend

### CategoryChart

Chart showing task distribution by category.

**Location:** `src/lib/features/analytics/components/CategoryChart.svelte`

**Usage:**
```svelte
<script>
  import { CategoryChart } from '$lib/features/analytics';
</script>

<CategoryChart 
  data={categoryData}
/>
```

**Props:**
- `data` (CategoryData): Category distribution data

**Features:**
- Doughnut chart
- Category colors
- Task count display
- Interactive segments

### OverdueTasksSummary

Summary component for overdue tasks.

**Location:** `src/lib/features/analytics/components/OverdueTasksSummary.svelte`

**Usage:**
```svelte
<script>
  import { OverdueTasksSummary } from '$lib/features/analytics';
</script>

<OverdueTasksSummary 
  overdueTasks={overdueTasks}
/>
```

**Props:**
- `overdueTasks` (Task[]): Array of overdue tasks

**Features:**
- Overdue task count
- Urgency indicators
- Quick action buttons
- Grouping by priority

### ProductivityInsights

Component displaying productivity insights and recommendations.

**Location:** `src/lib/features/analytics/components/ProductivityInsights.svelte`

**Usage:**
```svelte
<script>
  import { ProductivityInsights } from '$lib/features/analytics';
</script>

<ProductivityInsights 
  insights={productivityData}
/>
```

**Props:**
- `insights` (InsightsData): Productivity insights data

**Features:**
- Productivity score
- Trend analysis
- Recommendations
- Achievement badges

### TimePeriodSelector

Selector for analytics time periods.

**Location:** `src/lib/features/analytics/components/TimePeriodSelector.svelte`

**Usage:**
```svelte
<script>
  import { TimePeriodSelector } from '$lib/features/analytics';
  
  let selectedPeriod = $state('week');
</script>

<TimePeriodSelector 
  selectedPeriod={selectedPeriod}
  onPeriodChange={(period) => selectedPeriod = period}
/>
```

**Props:**
- `selectedPeriod` (string): Currently selected period
- `onPeriodChange` (function): Callback when period changes

**Features:**
- Period options (day, week, month, year)
- Custom date range
- Quick select buttons
- Date picker integration

### Analytics Skeleton Components

Loading skeletons for analytics components.

**Locations:**
- `src/lib/features/analytics/components/ChartSkeleton.svelte`
- `src/lib/features/analytics/components/StatisticsCardsSkeleton.svelte`

**Usage:**
```svelte
<script>
  import { ChartSkeleton, StatisticsCardsSkeleton } from '$lib/features/analytics';
</script>

<StatisticsCardsSkeleton />
<ChartSkeleton />
```

**Features:**
- Matches component layouts
- Smooth loading animations
- Accessible placeholders

## Component Usage Best Practices

### 1. Component Organization

- Use feature-based organization for business logic components
- Use shared components for reusable UI primitives
- Keep components focused and single-purpose
- Follow the established directory structure

### 2. Props and State Management

- Use Svelte 5 runes (`$state`, `$derived`, `$props`) for reactivity
- Type all props with TypeScript interfaces
- Provide default values for optional props
- Use `$bindable` for two-way binding when appropriate

### 3. Error Handling

- Implement proper error boundaries
- Use toast notifications for user feedback
- Provide loading states for async operations
- Handle edge cases (empty states, errors)

### 4. Accessibility

- Use semantic HTML elements
- Provide ARIA labels where needed
- Ensure keyboard navigation support
- Use proper color contrast
- Include loading indicators

### 5. Performance

- Use `$derived` for computed values
- Implement debouncing for search inputs
- Lazy load components when appropriate
- Use skeleton screens for better perceived performance

### 6. Testing

- Write unit tests for component logic
- Test component props and state
- Test user interactions
- Use test fixtures for complex data
- Mock external dependencies

### 7. Styling

- Use Tailwind CSS utility classes
- Follow mobile-first responsive design
- Use consistent spacing and colors
- Implement proper hover and focus states
- Ensure cross-browser compatibility

### 8. Data Flow

- Use stores for shared state management
- Pass data down through props
- Emit events up through callbacks
- Keep components stateless when possible
- Use server-side loading for initial data

### 9. Form Handling

- Use Zod schemas for validation
- Provide clear error messages
- Implement proper form reset
- Handle submission states
- Validate on blur and submit

### 10. Internationalization

- Keep text external to components when possible
- Use locale-aware date formatting
- Consider RTL language support
- Provide translatable error messages

## Component Examples

### Complete Task Management Page Example

```svelte
<script lang="ts">
  import { TaskList, TaskFilters, TaskSearch, TaskForm } from '$lib/features/tasks';
  import { taskStore } from '$lib/features/tasks/stores/task.store';
  import { LoadingSpinner, EmptyState } from '$lib/shared/components';
  import { toastStore } from '$lib/shared/stores';
  
  let showForm = $state(false);
  let editingTask = $state(null);
  
  async function handleCreateTask(data) {
    try {
      await taskStore.createTask(data);
      toastStore.success('Task created successfully');
      showForm = false;
    } catch (error) {
      toastStore.error('Failed to create task');
    }
  }
  
  async function handleUpdateTask(data) {
    try {
      await taskStore.updateTask(editingTask.id, data);
      toastStore.success('Task updated successfully');
      editingTask = null;
      showForm = false;
    } catch (error) {
      toastStore.error('Failed to update task');
    }
  }
  
  async function handleDeleteTask(task) {
    try {
      await taskStore.deleteTask(task.id);
      toastStore.success('Task deleted successfully');
    } catch (error) {
      toastStore.error('Failed to delete task');
    }
  }
</script>

<div class="container mx-auto p-4">
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-2xl font-bold">Tasks</h1>
    <button 
      on:click={() => showForm = true}
      class="bg-blue-500 text-white px-4 py-2 rounded"
    >
      Add Task
    </button>
  </div>
  
  {#if showForm}
    <TaskForm 
      mode={editingTask ? 'edit' : 'create'}
      initialData={editingTask}
      onSubmit={editingTask ? handleUpdateTask : handleCreateTask}
      onCancel={() => {
        showForm = false;
        editingTask = null;
      }}
    />
  {/if}
  
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
    <TaskSearch />
    <TaskFilters />
  </div>
  
  {#if taskStore.state.isLoading}
    <LoadingSpinner />
  {:else if taskStore.state.tasks.length === 0}
    <EmptyState 
      title="No tasks found"
      description="Create your first task to get started"
      icon="📝"
    />
  {:else}
    <TaskList 
      tasks={taskStore.state.tasks}
      onEditTask={(task) => {
        editingTask = task;
        showForm = true;
      }}
      onDeleteTask={handleDeleteTask}
    />
  {/if}
</div>
```

## Additional Resources

- [Svelte 5 Documentation](https://svelte.dev/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Zod Validation](https://zod.dev/)
- [Chart.js Documentation](https://www.chartjs.org/docs/)
- [AGENTS_FrontEnd.md](./sveltekit-taskmanager-frontend/AGENTS_FrontEnd.md) - Development guidelines
