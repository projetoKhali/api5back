import pandas as pd
from extraction_data import *
from datetime import datetime

def create_dim_datetime():
    """Cria um registro na dimensão datetime com o momento atual."""
    data_atual = datetime.now()
    dados = {
        'date': data_atual.date(),
        'year': data_atual.year,
        'month': data_atual.month,
        'weekday': data_atual.weekday(),
        'day': data_atual.day,
        'hour': data_atual.hour,
        'minute': data_atual.minute,
        'second': data_atual.second
    }
    df_dim_datetime = pd.DataFrame([dados])
    df_dim_datetime['date'] = pd.to_datetime(df_dim_datetime['date']).dt.date
    df_dim_datetime['year'] = df_dim_datetime['year'].astype('int64')
    df_dim_datetime['month'] = df_dim_datetime['month'].astype('int64')
    df_dim_datetime['weekday'] = df_dim_datetime['weekday'].astype('int64')
    df_dim_datetime['day'] = df_dim_datetime['day'].astype('int64')
    df_dim_datetime['hour'] = df_dim_datetime['hour'].astype('int64')
    df_dim_datetime['minute'] = df_dim_datetime['minute'].astype('int64')
    df_dim_datetime['second'] = df_dim_datetime['second'].astype('int64')

    return df_dim_datetime

df_dim_datetime = create_dim_datetime()

df_dim_user = df_user.rename(columns={
    'usr_id': 'id',
    'usr_name': 'name',
    'usr_ocupation': 'occupation'
})[['id', 'name', 'occupation']]

df_dim_process = df_process.rename(columns={
    'pc_id': 'id',
    'pc_title': 'title',
    'pc_initial_date': 'initial_date',
    'pc_finish_date': 'finish_date',
    'pc_status': 'status',
    'usr_id': 'dim_usr_id',
    'pc_description': 'description'
})[['id', 'title', 'initial_date', 
    'finish_date', 'status', 'dim_usr_id', 'description']]

df_dim_vacancy = df_vacancy.rename(columns={
    'vc_id': 'id',
    'vc_title': 'title',
    'vc_num_positions' : 'num_positions',
    'vc_status': 'status',
    'vc_location' : 'location',
    'usr_id' : 'dim_usr_id',
    'vc_opening_date' : 'opening_date',
    'vc_closing_date' : 'closing_date'
})[['id', 'title', 'num_positions', 
    'status', 'location', 'dim_usr_id', 'opening_date', 'closing_date']]

def create_fact_table(df_process, df_vacancy, df_vacancy_candidate, df_interview, df_feedback, df_hiring):
    """Cria a tabela fato com todas as métricas e dimensões necessárias."""
    
    fact_table = pd.merge(df_process, df_vacancy, left_on='pc_id', right_on='pc_id', suffixes=('_process', '_vacancy'))
    
    candidates_per_vacancy = df_vacancy_candidate.groupby('vc_id').size()
    fact_table['met_total_candidates_applied'] = fact_table['vc_id'].map(candidates_per_vacancy).fillna(0).astype(int)
    
    interviews_per_vacancy = df_interview.groupby('vc_id').size()
    fact_table['met_total_candidates_interviewed'] = fact_table['vc_id'].map(interviews_per_vacancy).fillna(0).astype(int)
    
    hires_per_vacancy = df_hiring.groupby('vc_id').size()
    fact_table['met_total_candidates_hired'] = fact_table['vc_id'].map(hires_per_vacancy).fillna(0).astype(int)
    
    fact_table['met_sum_duration_hiring_proces'] = (pd.to_datetime(fact_table['pc_finish_date']) - pd.to_datetime(fact_table['pc_initial_date'])).dt.days.fillna(0).astype(int)
    
    salary_sum = df_hiring.groupby('vc_id')['hr_initial_salary'].sum()
    fact_table['met_sum_salary_initial'] = fact_table['vc_id'].map(salary_sum).fillna(0).astype(int)
    
    feedback_counts = df_feedback.groupby(['vc_id', 'fd_type']).size().unstack(fill_value=0)
    feedback_counts = feedback_counts.reindex(columns=[1, 2, 3], fill_value=0)  # Ensure all columns exist
    fact_table = fact_table.merge(feedback_counts, left_on='vc_id', right_index=True, how='left')
    
    fact_table = fact_table.rename(columns={
        'pc_id': 'dim_process_id',
        'vc_id': 'dim_vacancy_id',
        'usr_id_process': 'dim_user_id',
        1: 'met_total_feedback_positive',
        2: 'met_total_negative',
        3: 'met_total_neutral'
    })
    
    fact_table = fact_table.reset_index(drop=True)
    fact_table['id'] = fact_table.index + 1
    
    columns = [
        'id',
        'met_total_candidates_applied',
        'met_total_candidates_interviewed',
        'met_total_candidates_hired',
        'met_sum_duration_hiring_proces',
        'met_sum_salary_initial',
        'met_total_feedback_positive',
        'met_total_neutral',
        'met_total_negative',
        'dim_process_id',
        'dim_vacancy_id',
        'dim_user_id'
    ]
    
    fact_table = fact_table[columns]
    
    for col in fact_table.columns:
        fact_table[col] = fact_table[col].fillna(0).astype(int)
    
    return fact_table

fact_hiring_process = create_fact_table(df_process, 
        df_vacancy, 
        df_vacancy_candidate, 
        df_interview, 
        df_feedback, 
        df_hiring)

if __name__ == "__main__":
    fact_hiring_process = create_fact_table(
        df_process, 
        df_vacancy, 
        df_vacancy_candidate, 
        df_interview, 
        df_feedback, 
        df_hiring
    )
    
    print("Dimensão DateTime:")
    print(df_dim_datetime.head())

    print("Dimensão User:")
    print(df_dim_user.head())

    print("\nDimensão Process:")
    print(df_dim_process.head())

    print("\nDimensão Vacancy:")
    print(df_dim_vacancy.head())
    
    print("\nFact Hiring Process:")
    print(fact_hiring_process.head())