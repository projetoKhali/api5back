import pandas as pd
from extraction_data import *
from datetime import datetime

def transformar_dim_datetime(data):
    """Converte uma data para um dicionário de atributos para a dimensão datetime."""
    return {
        'dim_datetime_id': int(data.timestamp()),
        'dim_datetime_date': data.date(),
        'dim_datetime_year': data.year,
        'dim_datetime_month': data.month,
        'dim_datetime_weekday': data.weekday(),
        'dim_datetime_day': data.day,
        'dim_datetime_hour': data.hour,
        'dim_datetime_minute': data.minute,
        'dim_datetime_second': data.second
    }

def create_dim_datetime(num_records):
    """Cria registros na dimensão datetime baseados no momento atual."""
    data_atual = datetime.now()
    return pd.DataFrame([transformar_dim_datetime(data_atual) for _ in range(num_records)])

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
    
    # Criando a base da tabela fato
    fact_table = pd.merge(df_process, df_vacancy, left_on='pc_id', right_on='pc_id', suffixes=('_process', '_vacancy'))
    
    # Calculando métricas com tratamento de NA
    # Total de candidatos que se candidataram
    candidates_per_vacancy = df_vacancy_candidate.groupby('vc_id').size()
    fact_table['met_total_candidates_applied'] = fact_table['vc_id'].map(candidates_per_vacancy).fillna(0).astype(int)
    
    # Total de candidatos entrevistados
    interviews_per_vacancy = df_interview.groupby('vc_id').size()
    fact_table['met_total_candidates_interviewed'] = fact_table['vc_id'].map(interviews_per_vacancy).fillna(0).astype(int)
    
    # Total de candidatos contratados
    hires_per_vacancy = df_hiring.groupby('vc_id').size()
    fact_table['met_total_candidates_hired'] = fact_table['vc_id'].map(hires_per_vacancy).fillna(0).astype(int)
    
    # Soma da duração do processo de contratação (em dias)
    fact_table['met_sum_duration_hiring_process'] = (
        pd.to_datetime(fact_table['pc_finish_date']) - pd.to_datetime(fact_table['pc_initial_date'])
    ).dt.days.fillna(0).astype(int)
    
    # Soma dos salários iniciais
    salary_sum = df_hiring.groupby('vc_id')['hr_initial_salary'].sum()
    fact_table['met_sum_salary_initial'] = fact_table['vc_id'].map(salary_sum).fillna(0).astype(int)
    
    # Contagem de feedbacks por tipo
    feedback_counts = df_feedback.groupby(['vc_id', 'fd_type']).size().unstack(fill_value=0)
    feedback_counts = feedback_counts.reindex(columns=[1, 2, 3], fill_value=0)
    
    # Ensure all columns exist
    fact_table = fact_table.merge(feedback_counts, left_on='vc_id', right_index=True, how='left')
    
    # Criar dimensão datetime para cada registro
    df_dim_datetime = create_dim_datetime(len(fact_table))
    fact_table['dim_date_id'] = df_dim_datetime['dim_datetime_id']
    
    # Renomear colunas
    fact_table = fact_table.rename(columns={
        'pc_id': 'dim_process_id',
        'vc_id': 'dim_vacancy_id',
        'usr_id_process': 'dim_user_id',
        'fac_id': 'id',
        'dim_datetime_id' : 'dim_date_id',
        1: 'met_total_feedback_positive',
        2: 'met_total_neutral',
        3: 'met_total_negative'
    })
    
    # Criar ID e reorganizar colunas na ordem correta
    fact_table = fact_table.reset_index(drop=True)
    fact_table['fact_id'] = fact_table.index + 1
    fact_table = fact_table.rename(columns={'fact_id': 'id'})
    
    # Definir a ordem final das colunas
    columns = [
        'id',
        'met_total_candidates_applied',
        'met_total_candidates_interviewed',
        'met_total_candidates_hired',
        'met_sum_duration_hiring_process',
        'met_sum_salary_initial',
        'met_total_feedback_positive',
        'met_total_neutral',
        'met_total_negative',
        'dim_process_id',
        'dim_vacancy_id',
        'dim_user_id',
        'dim_date_id'
    ]
    
    # Reordenar as colunas e converter para int
    fact_table = fact_table[columns]
    
    # Final NA check and conversion to int
    for col in fact_table.columns:
        fact_table[col] = fact_table[col].fillna(0).astype(int)
        
    return fact_table, df_dim_datetime

if __name__ == "__main__":
    # Criar as tabelas
    fact_hiring_process, dim_datetime = create_fact_table(
        df_process, 
        df_vacancy, 
        df_vacancy_candidate, 
        df_interview, 
        df_feedback, 
        df_hiring
    )
    
    print("Dimensão DateTime:")
    print(dim_datetime.head())

    print("Dimensão User:")
    print(df_dim_user.head())

    print("\nDimensão Process:")
    print(df_dim_process.head())

    print("\nDimensão Vacancy:")
    print(df_dim_vacancy.head())
    
    print("\nFact Hiring Process:")
    print(fact_hiring_process.head())