<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.title" placeholder="Title" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.importance" placeholder="Imp" clearable style="width: 90px" class="filter-item">
        <el-option v-for="item in importanceOptions" :key="item" :label="item" :value="item" />
      </el-select>
      <el-select v-model="listQuery.type" placeholder="Type" clearable class="filter-item" style="width: 130px">
        <el-option v-for="item in calendarTypeOptions" :key="item.key" :label="item.display_name+'('+item.key+')'" :value="item.key" />
      </el-select>
      <el-select v-model="listQuery.sort" style="width: 140px" class="filter-item" @change="handleFilter">
        <el-option v-for="item in sortOptions" :key="item.key" :label="item.label" :value="item.key" />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        Search
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        Add
      </el-button>
      <el-button v-waves :loading="downloadLoading" class="filter-item" type="primary" icon="el-icon-download" @click="handleDownload">
        Export
      </el-button>
      <el-checkbox v-model="showReviewer" class="filter-item" style="margin-left:15px;" @change="tableKey=tableKey+1">
        reviewer
      </el-checkbox>
    </div>
    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column width="120px" align="center" label="Name">
        <template slot-scope="scope">
          <span>{{ scope.row.category_name }}</span>
        </template>
      </el-table-column>
    <el-table-column width="200px" align="center" label="Create-Date">
        <template slot-scope="scope">
          <span>{{ scope.row.created_at  }}</span>
        </template>
      </el-table-column>
      <el-table-column width="200px" align="center" label="Update-Date">
        <template slot-scope="scope">
          <span>{{ scope.row.updated_at  }}</span>
        </template>
      </el-table-column>



      <el-table-column width="120px" align="center" label="Tags">
        <template slot-scope="scope">
          <span>{{ scope.row.category_tags }}</span>
        </template>
      </el-table-column>
      <el-table-column width="120px" align="center" label="Links">
        <template slot-scope="scope">
          <span>{{ scope.row.category_link }}</span>
        </template>
      </el-table-column>
      <el-table-column   label="Info">
        <template slot-scope="scope">
          <span>{{ scope.row.category_info }}</span>
        </template>
      </el-table-column>

      <!--


              <el-table-column class-name="status-col" label="Status" width="110">
                <template slot-scope="{row}">
                  <el-tag :type="row.status | statusFilter">
                    {{ row.status }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column min-width="300px" label="Title">
                <template slot-scope="{row}">
                  <router-link :to="'/example/edit/'+row.id" class="link-type">
                    <span>{{ row.title }}</span>
                  </router-link>
                </template>
              </el-table-column>
-->
              <el-table-column align="right" label="Actions"  width="120px">
                <template slot-scope="scope">
                  <router-link :to="'/resource/edit/'+scope.row.id">
                    <el-button type="primary" size="small" icon="el-icon-edit">
                      Edit
                    </el-button>
                  </router-link>
                </template>
              </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </div>
</template>

<script>
import { getCategory } from '@/api/resources-category'
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination


export default {
  name: 'ArticleList',
  components: { Pagination },
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'info',
        deleted: 'danger'
      }
      return statusMap[status]
    }
  },
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20
      }
    }
  },
  created() {
    this.getServerCategory()
  },
  methods: {
    getServerCategory() {
      this.listLoading = true
      getCategory(this.listQuery).then(response => {

        this.list = response.spec
        this.total = response.total
        this.listLoading = false
      })
    }
  }
}
</script>

<style scoped>
.edit-input {
  padding-right: 100px;
}
.cancel-btn {
  position: absolute;
  right: 15px;
  top: 10px;
}
</style>
